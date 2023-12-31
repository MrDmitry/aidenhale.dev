## Quick dip into Yocto

If you're not familiar with [Yocto project](https://www.yoctoproject.org/software-overview/), it's a powerhouse for
customized Linux systems fueled by the `bitbake` tool for the build management.

Yocto is based on a Layer system, where distribution is modularly configured using the software from the layers.
Usually the layers contain recipes to build components from source, but in cases of proprietary software (e.g. drivers)
you may work with a pre-built archive for some generic compatible architecture.

In my current project, we know all architectures and distributions our internal clients will integrate us with.
Source-based distribution was considered the last resort, as managing security access was known to be tricky:
* there are over a thousand engineers under different organizations that will work with our delivery
* developers may change teams and that would not be broadcast to a broader organization
* we'd need to establish a centralized access control over the cascading user groups (not only our direct users, but
their users as well), or allow our direct users to redistribute our delivery, neither of which sounds appealing

From performance standpoint, building from source would simply take too long and we couldn't efficiently share our
build caches with them, so we decided to prototype a
[pre-built](https://docs.yoctoproject.org/dev-manual/prebuilt-libraries.html) layer approach.

## Pre-built layer

I've worked with pre-built _libraries_ for Yocto before, mostly as part of proprietary vendor delivery, but it was
never a full layer. Yet it sounded very doable - after all, what is a layer if not a collection of libraries?

The scope of when we should stop packaging was not fully understood. We were sure that we don't want to package
low-level recipes (e.g. `glibc` or `zlib`), and we also were sure that we definitely want to redistribute our and some
vendor projects. But where was the line? We saw that vendor deployed custom patches for some recipes to enable some
behavior, that could brake ABI (or at least someone would have to monitor that), and we also saw vendor overriding some
package versions (usually downgrading) to ensure that it works as expected.

With this [paradox of the heap](https://en.wikipedia.org/wiki/Sorites_paradox) on our hands, we decided that it's
better to just start with something and then iterate. We knew some hard dependencies and license restrictions, so we
had enough direction to get us started. We also knew that there are red flags to consider any third-party project for
redistribution - patched being deployed, version being explicitly set and so on.

## Investigation and prototyping

Even though we knew all intended target systems, we wanted to have a generic system that would require minimal
maintenance regardless of the target environment. Unfortunately, we could not piggy-back off any existing package
manager supported by Yocto, as we knew that our end users will use different package managers. We could utilize the
generated **metadata** for the packages, namely things that `bitbake` generates as input for the package managers:
* package name and version
* runtime dependencies
* installed files

We also have access to the staging directory of each recipe, it contains all installed files.

Out idea was to package the staging directory as a tarball and generate a recipe based on the metadata, we could
recreate the exact staging environment prior to the packaging stage. Once that is done, we'd let the end-user
configuration to manage how recipes are packaged for their Linux distribution - effectively we've recreated the staging
directory and all needed instructions on how to use it.

The first prototype was fully manual, just to prove that approach was viable. We selected a short list of recipes to
try and capture the most common use cases:
* `-dev` only package
* common library
* application that uses common library

The second prototype was written in python and worked directly with the temporary build/staging directories. To
generate the recipes we used the original recipe file, but we filtered out all unexpected entries.

## Iterating on the prototype

The prototype was very simple - it only needed a path to the layer directory where it would parse the recipes from.
However it was unable to distinguish between build-time dependencies (e.g. host tools, static libraries or header-only
dependencies) and interface dependencies (e.g. header-only dependencies referred by the package). It doesn't sound
critical, but our users were confused as to why would certain tools and libraries be needed (and if they were needed,
where should they get them from), since they expected that our binary delivery is mostly self-contained. While
investigating we discovered that vast majority of those dependencies were no longer needed. They were a result of us
using the **original** recipe files and leaving the dependency information as-is.

This was not easily fixable, as even `bitbake` doesn't really know the true dependencies. Ultimately, it depended on
how the recipe maintainer organized their project and their recipes. There were some patterns that we noticed:
* host tools would only be needed to build the project and would not be required to use the generated binaries
* most header-only dependencies were not transitive - `boost` would be used in the project, but it would not be exposed
to the user via installed headers
* some domain-specific build-time dependencies were not transitive as well, as they were statically linked

We tried automating the detection, but we were always defeated by some edge case that was already present in our
delivery, for example there would be a manually maintained CMake package template `XXXConfig.cmake.in` that would
explicitly call `find_package(XYZ REQUIRED)`. This template was then used as a source for the installed CMake package
configuration. Nothing in the package configuration would tell us that this is not just a runtime dependency, but also
a build-time dependency. `bitbake` simply falls back to cascading `-dev` package dependencies because it doesn't
need to know any better, and it's rather cheap to generate hard links in the host environment.

Eventually we settled on adding filtering configuration to exclude arbitrary patterns from arbitrary fields. Looking
closely at `bitbake`-generated metadata also gave us better insights on dependencies, and we were able to create a
stable configuration, even though it required some minimal maintenance to ensure the filtering behaved as expected.

## Rewrite

This implementation was used for 3 release cycles and was then rewritten from scratch to utilize
[tinfoil](https://wiki.yoctoproject.org/wiki/TipsAndTricks/Tinfoil) that allowed us to properly process `bitbake`
recipes instead of parsing the generated metadata.

We knew that we no longer needed the original recipe, as `bitbake` metadata in combination with `tinfoil` allows us to
generate recipes for our prebuilts. We used [jinja](https://jinja.palletsprojects.com/en/2.11.x/) to template the
generated recipes. While switching, we fixed some bugs in our filtering logic - they were not obvious in the recipes
that utilized original recipe files, but switching to a single template uncovered some inconsistencies.

## Lessons learned

### Dependencies are hard

Initially we considered build-time dependencies to not be important for the prebuilts, after all as long as you can
extract it, no extra dependencies are needed. But that's not necessarily true for `-dev` packages, as they most
definitely are build-only dependencies. It'd be extremely complex for us to try and analyze build-time dependencies of
any `-dev` package that we generated - the dependency could be in a `XXXConfig.cmake` module, or in any other
configuration file. We chose to maintain a list of filters that are inherited from the layer to the recipe to allow
fine control over the `DEPENDS` parameter of the generated recipes.

### Sometimes flexibility is worse

We chose to keep a manifest of recipes that were allowed to be packaged for prebuilts. It was a bit annoying to
maintain, but overall new recipes are not added on a daily basis, so at worst the manifest will be updated the
following day, whenever automated tests start failing. Potentially we could organize our layers in some particular
format (e.g. utilizing some phony `.bbclass` to _tag_ a recipe for prebuilt generation), but this would be rather
intrusive and would make content management less obvious, as such changes could come from arbitrary `.bbappend` file
from any other layer.

### Users do the integration

We used a custom `.bbclass` to extract our "staging" tarball respecting the symbolic links and any other custom
deployment rules. This approach allowed our customers to override any paths if they needed, for example some of our
customers utilized `lib64` directory while others insisted on `lib` used even on 64-bit systems. Now they could override
the `.bbclass` with their own implementation to ensure that their rules are followed regardless of our tarball layout.

Preparing a full layer for pre-built integration is unviable - there always are build-only packages that should not be
exposed to the customers. Also there are external tools and libraries that we should not (or can not) distribute and
leave for our customer to resolve in their configuration. We can detect those and provide version-related information
together with our delivery.

Distributing [packagegroups](https://docs.yoctoproject.org/3.1.29/ref-manual/ref-classes.html?highlight=packagegroup#packagegroup-bbclass)
in combination with pre-built recipes improved the deployment for all our customers. We wanted to reuse the same
packagegroups in the generated pre-built layers, as we have in our development layers, so we added necessary logic to
also copy the packagegroups from our development layer to the pre-built layer.
