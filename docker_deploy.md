# In-container docker deployment

## Overview

This approach allows deploying software into a single docker image layer
while minimizing installer footprint via volume mounts.

Most of the time packages can be deployed through `Dockerfile` directly or, if
additional preparations are required, through multi-stage deployment.

But in some cases, usually with proprietary software, you only have a
self-installing wizard script.

Even if such script supports headless installation, it may have system-wide
side-effects making it close to impossible to multi-stage copy the necessary
installed artifacts.

## Background

A vendor provided several archives with inputs for several platform-specific
toolchains:
* `common-tools.tar.gz`
* `platform-A.tar.gz`
* `platform-B.tar.gz`

Platform-specific archives are installed on top of common tooling, already
suggesting some level of modularity to our installation.

Platform-specific archives can be deployed via multi-stage since they do not
have any side-effects, but common tools contain a self-installing wizard script
with some additional archives.

## Wizard analysis

Installer script performed a combination of `apt` installations (controlled by
the vendor) as well as extracting companion archives to a vendor location. Some
installation steps prompted for user input via `read` and did not support
`noninteractive` frontend configuration.

`apt` installations require `sudo` access (as per vendor implementation), while
other packages are installed to `~/.local`.

Maintaining vendor-controlled installation endpoint is not maintainable. Vendor
agreed to make the process more configurable for headless installation, but the
size of the installer was >500M that would be captured by an intermediate image
layer even if removed immediately after.

With no mount
