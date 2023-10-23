## Overview

This approach allows deploying software into a single docker image layer while minimizing installer footprint via
volume mounts.

Most of the time packages can be deployed through `Dockerfile` directly or, if additional preparations are required,
through multi-stage deployment.

But in some cases, usually with proprietary software, you only have a self-installing wizard script.

Even if such script supports headless installation, it may have system-wide side effects making it close to impossible
to multi-stage copy the necessary installed artifacts.

At the time, experimental feature of `BuildKit` allowed temporary bind-mounts via `RUN --mount` but due to corporate
policy we were not allowed to use `BuildKit` to ensure consistent `Dockerfile` syntax and usability.

## Background

A vendor provided several archives with inputs for several platform-specific toolchains:
* `common-tools.tar.gz`
* `platform-A.tar.gz`
* `platform-B.tar.gz`

Platform-specific archives are installed on top of common tooling, already suggesting some level of modularity to our
installation.

Platform-specific archives can be deployed via multi-stage since they do not have any side effects, but common tools
contain a self-installing wizard script with some additional archives.

## Baseline for future comparison

To evaluate how effective our attempts are, we gathered the initial metrics of `Dockerfile`-based deployment by calling
the wizard from an [expect](https://linux.die.net/man/1/expect) script.

Results looked something like this:
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

We know that there's at least **~30%** (653MB) of wasted space due to copying of the installer and its archives.

## Wizard analysis

Installer script performed a combination of `apt` installations (controlled by the vendor) as well as extracting
companion archives to a vendor location. Some installation steps prompted for user input via `read` and did not support
`noninteractive` frontend configuration.

`apt` installations require `sudo` access (as per vendor implementation), while other packages are installed to
`~/.local`.

Maintaining vendor-controlled installation endpoint is not maintainable. Vendor agreed to make the process more
configurable for headless installation, but the size of the installer was >500M that would be captured by an
intermediate image layer even if removed immediately after.

With no mounting capabilities, and with the `apt` side effects, the team decided to proceed with in-container
deployment with a subsequent `docker commit` of the resulting container.

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

Even though the main problem was resolved, there were runtime issues discovered, they were mostly related to how
non-portable the vendor software is:
* Developers use different UIDs (necessary UID shimming is done on the client side)
* Vendor software is installed (and usable) only for the UID at the time of installation
* Some additional configuration was required

An easy workaround of `chmod -R g=u` solved portability (ensuring that shimming adds the GID of the baked-in user) and,
coupled with additional configuration via another Dockerfile, it effectively doubled the image size:
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

This was expected, since image layers are additive and any file attribute change results in a copy of that file.

The only way to address this would be to perform any file modifications in the same layer the files are created. This
significantly increases the complexity of our installation scripts inviting re-evaluation of the whole approach.

### Re-evaluating the approach

Using `expect` script to watch the installation wizard proved relatively easy, [Tcl](https://en.wikipedia.org/wiki/Tcl)
language is rather easy to pick up.

However we discovered that we needed to perform not only post-install operations, but pre-install as well (but these
should not bleed into the final image).

Given the control we needed over the in-container temporary image reconfiguration, I decided to use `python` instead of
`bash`. This also gave us the flexibility to drop `expect` dependency from the target image and bring it to the host
environment - as [pexpect](https://pexpect.readthedocs.io/en/stable/) python module.

## Python automation

In addition to creating a wrapper for Docker container runtime, we decided to separate the deployment in 3 phases:
- `prepare` to process the input files (a few tarballs) and prepare the work environment for deployment
- `deploy` to perform deployment-related actions
- `test` to run basic sanity checks to ensure the installed tools are usable

First iteration had an overloaded user interface with ~10 parameters, most of which had default values and would rarely
be specified by the user. In the future iterations we moved some of these to be task-specific variables and simplified
the interface significantly, so the user only had to specify 1-3 parameters.

Platform-specific tasks were separated to independent python modules that are dynamically imported based on user input.

This in turn enabled straight-forward CI/CD automation of Docker image deployment. The automation job specified the
inputs to the scripts:
* Task to run
* Tag of base Docker image to be used
* Tag of target Docker image to be generated
* List of input files to be downloaded from file server(s)

Based on that input, python script would prepare the workspace, deploy the software, test the image and push it to the
registry for further validation.

## Lessons learned

Overall this was a very positive experience. In a sense we re-implemented the `RUN --mount` command based on our own
use cases, but I also learned a lot in the process. Using python instead of bash allowed us to contain the complexity
of the scripts so that any developer could pick up any part of the deployment process and maintain this script
ecosystem with maximum reusability. Would it be possible in bash? Yes. Should you do it in bash? No.

The first iteration of the tool was running great for about a year and was used to migrate through 4 major updates of
the vendor tooling. I ~~rewrote~~ refactored these scripts for better maintainability and improved CI/CD automation of
image deployment. This second iteration improved quality of life for state management of the in-container bash session:
* ephemeral credential, environment variable and proxy configuration usage based on host system configuration without
in-image contamination
* custom [context managers](https://docs.python.org/3/library/stdtypes.html#typecontextmanager) for managing proxy
settings as well as `sudo` scope
* added [argparse subparsers](https://docs.python.org/3/library/argparse.html#argparse.ArgumentParser.add_subparsers)
to allow per-task parameters instead of the catch-all interface of "all possible customizable settings"
