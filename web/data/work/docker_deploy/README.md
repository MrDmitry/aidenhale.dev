## Problem

We needed to efficiently deploy and distribute a rapidly evolving toolchain including a third-party vendor SDK,
tooling, and additional development tools and libraries. The resulting Docker image was intended for use by hundreds of
developers and in-cloud build agents, with an initial update cadence is estimated to be 2-3 times a week.

## Overview

While the additional tools and libraries could easily be installed via `Dockerfile` directly or via multi-stage, the
challenge lay in the third-party vendor SDK and tooling, which utilized a self-installing wizard script.

At the time, experimental feature of `BuildKit` allowed temporary bind-mounts via `RUN --mount` but due to corporate
policy we were not allowed to use `BuildKit` to ensure consistent `Dockerfile` syntax and usability.

We decided to experiment with in-container deployment via `docker run` with `docker commit` once ready. Such approach
allows deploying software into a single docker image layer while minimizing installer footprint via volume mounts.

## Background

Vendor SDK utilized independent archives for several platform-specific toolchains, suggesting modularity to the
deployment process:
* `common-tools.tar.gz`
* `platform-A.tar.gz`
* `platform-B.tar.gz`

Platform-specific archives were installable through a multi-stage deployment, but the common tools included a
self-installing wizard script with additional archives.

## Baseline for future comparison

To evaluate how effective our attempts are, we measured the naive `Dockerfile`-based deployment by calling the wizard
from an [expect](https://linux.die.net/man/1/expect) script, revealing approximately **30%** (653MB) of wasted space
due to the unnecessary copy of the installer and its archives:
```bash
$ docker images deploy
REPOSITORY   TAG        IMAGE ID       CREATED         SIZE
deploy       baseline   f9ab681fee29   3 seconds ago   1.99GB

$ docker history deploy:baseline
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
f9ab681fee29   5 seconds ago    /bin/sh -c cd /tmp/installer &&     ./expect…   1.21GB    
30d44ef75261   18 seconds ago   /bin/sh -c #(nop) COPY dir:ea460903c0553546a…   653MB     
cc72face2225   22 seconds ago   /bin/sh -c #(nop) WORKDIR /opt                  0B        
82bf53d5a14b   22 seconds ago   /bin/sh -c apt-get update && apt-get install…   54.6MB    
...
```

## Wizard analysis

Installer script combined `apt` installations (controlled by the vendor) with companion archives extractions to a
vendor location. Some installation steps prompted for user input via `read` and did not support `noninteractive`
frontend configuration. `apt` installations required `sudo` access (as per vendor implementation), while some packages
are installed to `~/.local`, creating non-portability issues for users with different UIDs among developers.

Maintaining vendor-controlled installation endpoint was considered unfeasible. Vendor agreed to make the process more
configurable for headless installation, but the size of the installer was over 500M and that would be captured by an
intermediate image layer even if removed immediately after.

We proceeded with the prototype, utilizing the `expect` script as an entry point to the `docker run` with the installer
directory mounted to the container runtime:
```bash
$ docker images deploy
REPOSITORY   TAG              IMAGE ID       CREATED         SIZE
deploy       prototype        f827f2d58bff   4 seconds ago   1.34GB
deploy       baseline         f9ab681fee29   7 minutes ago   1.99GB

$ docker history deploy:prototype
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
f827f2d58bff   11 seconds ago   /tmp/installer/expect_install                   1.21GB    install.sh
cc72face2225   7 minutes ago    /bin/sh -c #(nop) WORKDIR /opt                  0B        
82bf53d5a14b   7 minutes ago    /bin/sh -c apt-get update && apt-get install…   54.6MB    
...
```

This prototype removed the unwanted `COPY` of the installer package and we've recovered the wasted 653MB.

### Iterating on the prototype

The main problem was resolved, but we discovered some runtime issues that were related to the non-portability of the
vendor software:
* Developers use different UIDs (necessary UID shimming is done on the client side)
* The installed vendor software usable only for the installing user UID
* Additional configuration is required on the client side

An easy workaround of `chmod -R g=u` solved portability (ensuring that shimming adds the GID of the default user) and,
coupled with additional configuration via another `Dockerfile`, it effectively doubled the image size:
```bash
$ docker images deploy
REPOSITORY   TAG               IMAGE ID       CREATED          SIZE
deploy       prototype_fixed   d76b92a86c2f   7 seconds ago    2.54GB
deploy       prototype         aa64045658e6   16 seconds ago   1.34GB
deploy       baseline          2719cd017382   26 minutes ago   1.99GB

$ docker history deploy:prototype_fixed
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
d76b92a86c2f   18 seconds ago   /bin/sh -c chmod -R g=u /opt/vendor             1.21GB    
aa64045658e6   27 seconds ago   /tmp/installer/expect_install                   1.21GB    install.sh
cc72face2225   45 minutes ago   /bin/sh -c #(nop) WORKDIR /opt                  0B        
82bf53d5a14b   45 minutes ago   /bin/sh -c apt-get update && apt-get install…   54.6MB    
...
```

This was expected since image layers are additive and any file attribute changes result in a copy of that file. Even if
this is done on the client side, it would result in several minutes wasted every time there's a new image to be
configured.

The only way to address this would be to perform any file modifications in the same layer where the files are created.
This significantly increased the complexity of our installation scripts, prompting us to re-evaluate the approach.

### Re-evaluating the approach

Using `expect` script to watch the installation wizard proved relatively easy, [Tcl](https://en.wikipedia.org/wiki/Tcl)
language is rather easy to pick up. However we discovered that the script needs to perform both post- and pre-install
operations (and these should not bleed into the final image).

We decided to rewrite the prototype in `python` and replace the in-image dependency on the `expect` tool to host
dependency on the [pexpect](https://pexpect.readthedocs.io/en/stable/) python module.

## Python automation

In addition to creating a wrapper for the Docker container runtime, we separated the deployment to 3 phases:
- `prepare` to process the input files (a few tarballs) and prepare the host environment for deployment
- `deploy` to perform in-container deployment-related actions
- `test` to run basic sanity checks to ensure the installed tools are usable

First iteration had an overloaded user interface with ~10 parameters, most of which had default values and would rarely
be specified by the user. In the future iterations we moved some of these parameters to be task-specific and simplified
the interface significantly, so the user only had to specify 1-3 parameters.

Platform-specific tasks were separated to independent python modules that are dynamically imported based on user input.

This in turn enabled straight-forward CI/CD automation of Docker image deployment. The automation job specified the
inputs to the scripts:
* task to run
* tag of the base Docker image to be used
* tag of the target Docker image to be generated
* list of the input files to be downloaded from file server(s)

Based on that input, python script prepared the workspace, deployed the software, tested the generated image and pushed
it to the registry for further validation.

## Lessons learned

Overall this was a very positive experience. In a sense we re-implemented the `RUN --mount` command based on our own
use cases, and learned a lot in the process. Using `python` instead of `bash` allowed us to contain the complexity of
the scripts so that any developer could pick up any part of the deployment process and maintain this script ecosystem
with maximum reusability.

Would this be possible in bash? Yes. Should you do it in bash? No.

The first iteration of the tool was running great for about a year and was used to migrate through 4 major updates of
the vendor tooling. I ~~rewrote~~ refactored most of the tool for better maintainability and improved CI/CD automation
of image deployment. This second iteration improved quality of life for state management of the in-container bash
session:
* ephemeral credential, environment variable and proxy configuration usage mounted from the host
* custom [context managers](https://docs.python.org/3/library/stdtypes.html#typecontextmanager) for managing network
proxy settings as well as limiting the `sudo` scope
* switched to [argparse subparsers](https://docs.python.org/3/library/argparse.html#argparse.ArgumentParser.add_subparsers)
to allow per-task parameters instead of the catch-all interface of "all possible customizable settings"

The tool exceeded my expectations, as we were able to meet several goals that were initially considered unlikely:
* maximize layer reusability to minimize downloads on toolchain updates
* control how image layers are split to allow parallel downloads
* simplify migration activities for major toolchain upgrades
