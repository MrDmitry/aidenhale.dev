## Quick dip into Yocto

If you're not familiar with [Yocto project](https://www.yoctoproject.org/software-overview/), it's a powerhouse for
customized Linux systems fueled by the `bitbake` tool for the build management.

Yocto is based on Layer system, where distribution is modularly configured using the software from the layers. Usually
the layers contain recipes to build components from source, but in cases of proprietary software (e.g. drivers) you may
work with a pre-built archive for some generic compatible architecture.

In the case of my current project, we know all the architectures and distributions our internal clients will use us in.
Building from source would take too long and we couldn't efficiently share our build caches with them, so we decided
to prototype a [pre-built](https://docs.yoctoproject.org/dev-manual/prebuilt-libraries.html) layer approach.

## Pre-built layer

My previous experience with pre-built libraries for Yocto was with a limited subset of proprietary vendor software, it
never was a full layer. But it must be doable - after all, what is a layer if not a collection of libraries?

## Investigating and prototyping

Even though we know the target systems, I wanted to have a generic system that would require minimal maintenance
regardless of the target environment. We could not utilize any existing package manager provided by Yocto. However we
could utilize the generated metadata for the packages, namely things that `bitbake` generates as input for the package
managers:
* package name and version
* runtime dependencies
* installed files

We also have access to the recipe-wide staging directory, that contains installed files for all subprojects.

If we were to package the staging directory as a tarball and generate a recipe based on the metadata, we could recreate
the exact staging environment prior to the packaging stage, and then we'd let the end-user configuration to manage how
recipes are packaged for their Linux distribution.

Prototype implementation was done with python and did not utilize any `bitbake` API, it worked purely on the generated
files. To generate the recipes we used the original recipe files and filtering out all unexpected entries from it.

## Rewrite

This implementation was used for 2 or 3 release cycles and was then rewritten from scratch to utilize
[tinfoil](https://wiki.yoctoproject.org/wiki/TipsAndTricks/Tinfoil) that allowed us to properly process `bitbake`
recipes instead of parsing the generated metadata. We also used [jinja](https://jinja.palletsprojects.com/en/2.11.x/)
to template the generated recipes.

## Lessons learned

Initially we considered build-time dependencies to not be important for the pre-builts, after all as long as you can
extract it, no extra dependencies are needed. But that's not true for `-dev` packages for example, as they most
definitely are build-only dependencies. It'd be extremely complex for us to try and analyze build-time dependencies of
any `-dev` package that we generated - the dependency could be in a `FindXXX.cmake` module, or in any other
configuration file. We chose to maintain a list of filters that are inherited from the layer to the recipe to allow
fine control over the `DEPENDS` parameter of the generated recipes.

A custom `.bbclass` was used to extract our "staging" tarball respecting the symbolic links and any other custom
deployment rules. This approach allowed our customers to override any paths if they needed, for example some of our
customers utilized `lib64` directory while others insited on `lib` used for everything. Now they could override the
`.bbclass` with their own implementation to ensure that their rules are followed regardless of our tarball layout.

Preparing a full layer for pre-built integration is unviable - there always are build-only packages that should not be
exposed to the customers. Also there are external tools and libraries that we should not distribute and leave for our
customer to resolve in their configuration.

We chose to keep a manifest of recipes that were allowed to be packaged for pre-builts. It was a bit annoying to
maintain, but overall new recipes are not added on a daily basis, so at worst the manifest will be updated the
following day.

Distributing packagegroups in combination with pre-built recipes improved the deployment for all our customers. We
wanted to reuse the same packagegroups in the generated pre-built layers, as we have in our development layers, so we
added necessary logic to also copy the packagegroups from our development layer to the pre-built layer.
