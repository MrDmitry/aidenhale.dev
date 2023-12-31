## Monorepo is good

If your project is simple - monorepo is good.

If your project is enterprise - monorepo can be good.

If you don't care about coupling your projects - monorepo works.

If you enjoy resolving merge conflicts - monorepo is very good.

If you want to _feel_ monorepo dev workflow without actually going monorepo you can use tools for that:
* google's [repo tool](https://source.android.com/source/using-repo.html) and alternatives
* `git submodule` if you're feeling frisky

## Monorepo is bad

If you want to revert functional changes - monorepo is bad. But if you like resolving merge conflicts - it's good.

If you want to express "release artifacts" through your repository (e.g. via tags) - monorepo is pretty bad. It forces
a single digital heartbeat and doesn't leave much space for per-project release management.

## Monorepo is ugly

If your project has big PR throughput - monorepo is pretty ugly. But if it's enterprise you probably can afford PR
merge orchestration to ensure that you can get as many non-conflicting or easily-resolved PRs merged in a timely manner.

If your ecosystem doesn't allow in-tree/exported package resolution - monorepo is ugly. You need to manage dependencies
in your CI to properly cascade the checks for projects that depend on your change.

If you want to decouple "development" from "integration" - monorepo turns ugly. It's great if you can always integrate
when you develop, but the more distributed the contributors are, the larger your monorepo becomes, and the more common
dependencies get established, the harder it gets to iteratively improve legacy code base. Parts of your repository
turn into zombies waiting to be replaced and deprecated.

## What should I use?

If you have to ask, you're in trouble, may as well flip a coin.

And if you're just being snarky, you probably know what's best for your project, so "it depends".

I'm a multi-repo fan, so I'd advise starting with multi-repo and feeling it out - if you don't like it, you can easily
switch to monorepo.

I worked in a project with horrible multi-repo layout that had so much friction we ended up having week-long "code
freezes" where master branches were closed to contributions from non-admins. And that's how we had our quarterly
releases.

I worked in a project with horrible monorepo layout that had so much friction we also had "code freezes" and constantly
had problems with behavior differences between development and production.

Was it caused by multi-repo vs monorepo choice? Not really, it was caused by improper dependency management in both
cases. And I don't think choosing the other layout would've solved anything, the problem was elsewhere.

My biggest problem with monorepo is reverting changes. It's rare, but always messy - you bisect to find the breaking
change and then can't revert it, because it touched several projects and some of those had changes in the same areas so
a direct revert cannot apply and you have to involve other team(s) to integrate the revert because it's not just about
a merge conflict, it's also about fixing the surrounding logic because, well, some project moved faster than the others.

I like the flexibility that multi-repo brings where you can always revert to a particular revision of any project.
While development team is bisecting, any user of the offending project may revert to some older version. And since
integration with a newer version was likely handled through its own integration PR (in case of submodules) or via some
manifest update (in case of repo-like tools), it's easy to identify what needs to be reverted in user code to restore
some older functionality.

Multi-repo requires a bit more day-to-day maintenance to align the revisions of your dependencies (either submodules or
manifests), but you retain the independence of individual projects.

With monorepo you have to pay the upfront costs (and some maintenance overhead) to manage dependencies between your
projects. It may sound trivial, but you'll have to trade-off somewhere:
* should every dependency use the latest version? sure
* should every change in a dependency trigger a rebuild for all downstream projects? maybe
* should every PR run unit tests? sure
* should every PR run functional tests? sure
* should every PR run integration tests? uhhh maybe
* should every PR run end-to-end tests? maybe not

There's no silver bullet, it's always some trade-off:
* do you want to gate the development until your <insert test type> passes? sometimes, but probably not always
* do you want to gate the development of your dependency because of another component's flaky tests? not really
* do you want to never run <insert test type>? not really, you need to know when you fix downstream tests

## So how to CI with all of that?

I prefer to separate "development" and "integration" workflows.

"Development" is the shortest path for a change to be contributed. It could be associated with some feature, but may
not necessarily deliver the feature in full. Most of the time it's some incremental improvement or bug fixing.
`Required` checks usually reflect the "Development" quality gates. Metrics would be different from project to project,
but in general we want to gate the development only if some minimal required testing fails.

"Integration" is a longer path for a change, usually coupled with some cross-project coordination. It could be an
upgrade of some dependency, it could be a new feature rollout, it could be some internal milestone for the project.
Ideally we want our CI to know "when to run extended checks", but at the very least we want to be able to trigger
`Optional` extended testing at developer/reviewer discretion.

Multi-repo allows us to have separation of concerns within our CI - all orchestration is handled on repository level
and it's up to the repository owners to decide how to establish quality gates for their project.

Monorepo requires additional orchestration (potentially cascading) for the contents of the monorepo. If we don't want
to establish per-project quality gates, we may utilize some common set of quality gates for the whole thing. But simply
because it's a monorepo, we may discover the dark side of CI - the day-to-day Continuous Integration. Not only do you
need to stay up-to-date with the moving base branch, you also often need to perform _integration_ while _developing_.
There are approaches to evolve interfaces in a backwards-compatible fashion, but at the end of the day any change to a
monorepo is expected to contribute to a "stable and sane configuration", which is only possible if integration gates
pass. Do you have to run end-to-end testing for any change? No, but you have to be honest that your master branch may
not necessarily be healthy at any point in time.

If your ecosystem allows you to keep using some "older version published via a package manager" you don't immediately
bump into this problem, but you will have to establish some process of "using latest distributed by monorepo" and you
will have to figure out a way to trigger your end-to-end testing using your latest (potentially not integrated)
dependency.

Does multi-repo solve this? Not really, but it allows you to separate the development from integration in a cleaner
fashion.

Do you always need that? No, probably not.

Can you find a way to achieve the same with monorepo? Sure.

Can you figure it out when it becomes a problem? Sure.

You will still run your extended tests on a separate cadence (maybe even a rolling test for whatever got merged
between the test runs), but for me it's easier to reason about "development" and "integration" as separate activities,
and multi-repo just maps nicely to that.

## Automated integration

This is a tangent, but regardless of mono- or multi-repo setup, you can automate your integration to "always look at
latest" and even update your dependencies automatically if you trust your tests enough. The difference to me is the
level of autonomy given to any sub-project.

In a monorepo you intend to have a sane configuration, but there's always some trade-off to enable rapid development.

In a multi-repo you manage your "integration" as a separate unit. It could be a manifest that lists a "known sane and
stable configuration of projects", or it could be a separate repository that represents integration of known
submodule revisions.

It it pretty? No. Especially the submodules.

Is it easy to manage? Manually - annoying; automated - same as monorepo.

But it's much easier to reason about "integration" when it is separate from "development". And I'm not advocating for
"manual integration" or "big bang integration" or any other moniker of "future me will handle it", I'm only separating
integration onto its own timeline, when compared with a monorepo it will barely affect the "feature delivery" timeline.
