## Background

I was tasked with designing a Continuous Integration (CI) system for a middleware project. The scale of the project was
bigger than what I used to work with, so I decided to be a bit more thorough in my design:
* Embedded stack targeting at least 3 different operating systems
* Project consists of up to 100 individual components
* Over 300 developers contributing directly to the project
* Over 1000 developers consuming the middleware and integrating it into their toolchains

Unsurprisingly, I was not able to find much information online for considerations and overall approach to designing
automation for a project of this scale. Most of the articles I found were authored by some CI service provider pitching
some high-level _"What is CI anyway?"_ article as well as selling their services or solutions. Such articles would be a
mix of technical solutions (how to configure certain types of pipelines) and some design considerations (what kinds of
automation gates to create), but there is very little on the **values** that automation should aim to bring to the
project, or even general considerations when designing your own CI system.

There seems to be some consensus over how CI should look like, and more specifically how it should scale. But I was not
able to find anything on how to approach the design from scratch. So I decided to dive deeper and come up with such
an approach.

Throughout the design and implementation of this CI system I discovered that _automation as a tool_ constantly falls
victim to the [XY problem](https://en.wikipedia.org/wiki/XY_problem) and, regardless of the underlying problem,
stakeholders instinctively jump to _"this should be a part of automation"_ even when the actual problem may not be well
understood. This provides a false sense of security that the problem is not only understood, but also almost solved.

In this article I do not attempt to answer the _"how to automate Y"_ question, but instead _"how should automation be
approached"_.

## Goal

I started with the goal of designing a CI system, but quickly realized that CI was just a problem domain, while the tool
was **automation**. So I decided to explore automation through the lens of a CI system, with a final goal of identifying
a set of core values that form a foundation for any automation-related design work.

Main audience for this article are infrastructure engineers and their managers. I believe it will also be helpful to any
software engineer who wants to reflect on values and core principles automation revolves around, as it often means
different things to different people, often becoming a catch-all for _"I don't need to do X locally, automation will
take care of it."_

## Table of contents

- [What is automation?](#what-is-automation)
- [Core values](#core-values)
- [GOAT design framework](#goat-design-framework)
- [Putting the GOAT into action](#putting-the-goat-into-action)
- [Eliminating automation SMEs](#eliminating-automation-smes)
- [Healthy ecosystem](#healthy-ecosystem)
- [Achievements](#achievements)
- [Lessons learned](#lessons-learned)

## What is automation?

Automation is a set of tools and pipelines that transforms repetitive manual actions to a headless scripted execution.
By virtue of being headless, automation ensures that human factor is limited to the initial phase of capturing the
desired process into the headless scripted format (and any subsequent modifications to it).

Automation requires a _good enough_ understanding of the process-under-automation, meaning that automation has to
piggyback on the existing tools and processes developers use, and that it should not attempt to reimplement tools. This
ensures that there is minimum to no automation-specific logic, and any automated process is already known to a subset of
engineers on the project.

Automation-specific logic has to be captured as configuration of the automation service (e.g. Jenkins, TeamCity, GitHub
Actions, etc), but rarely as standalone tools. When standalone tools have to be developed, their usage should be
promoted among the relevant engineers.

Any automated process has to have a specific goal as well as some known limitations. In CI systems we could say that a
successful _PR build_ proves that a repository is buildable, but not necessarily that it can be integrated into the
target system, or that it behaves as intended.

As an example, goals for a PR build could be framed as:
* Ensure base branch is stable (can be built) for the target platform
* Produce a set of artifacts that can be tested

With these goals we also have some limitations that should be acknowledged:
* Depending on build configuration, it could represent pre-merge or post-merge repository state
* Build result does not represent if produced artifacts are usable

## Core values

Here are the values I propose to be the foundation for any automation.

### Scoped goal

The goal for any automated process has to be specific and measurable, you should also try to scope it to some smallest
meaningful "gate". It's much easier to **compose** an automated _system_ from individual _parts_ than to **break up** a
_monolith_.

Depending on project layout and complexity, you may not have "integration" process at all, and all your goals may fit in
a build-and-test job, but still consider those gates to be separate - build proves that it can be built, test proves
that it behaves a certain way - no need to tangle the two. But if a project is rather small/simple - might as well
combine them.

Once you scope the goals for automation, you can treat these goals as "requirements" for your system. You may not have
expressed them as such, but technically they are describing specifically what your automation achieves. Might as well
use them as such.

### Acknowledged limitations

There are always edge cases that automation smooths over, but the limitations stay there, even if they are not as
obvious.

Consider automating a release process for several inter-connected projects. You may want to tag the release candidate
_before_ you attempt your release build, or you may want to tag the release candidate _after_ all critical gates are
passed. Whatever you choose poses a limitation to your release process - either a tagged revision can be unstable or
unsuitable for release, or you may not get a tagged revision because of some failures.

Acknowledging limitations allows you to simplify your triage or analysis of automation behavior, they can serve as a
checklist of _"did X happen?"_ items that will guide you through any undesired outcome.

### Iteration on goals and limitations

Every time automation fails (either the pipelines are failing, or the outputs are incorrect/unexpected) revisit the
goals and limitations and check if they actually describe your current system. It is possible that stakeholder assumed
that automation did certain things and goals should be extended. Likewise, it is possible that there were unexpected
limitations resulting in undesired outputs.

The better you understand (and describe) your automation, the easier it is to maintain it. Be clear with all
stakeholders what specifically is automated, and what is the extent of such automation.

## GOAT design framework

I wanted a silly acronym for this framework, so I went with GOAT - Goal-Oriented Automation Toolkit.

To create a GOATed automation you need to address two main parts:
* How to **approach a problem** with automation in mind?
* How to **design** an individual automated process?

### Approaching the problem

Main question is simple: **do we already utilize some manual solution for this problem?**

If the answer is **NO**, you should plan for an additional investigatory phase to understand the problem at hand. The
efforts always depend on the problem, but it is not viable to attempt to design automation for a _new problem_.

If the answer is **YES**, or if you're already planning for hands-on investigation, make sure that _inputs_ are known
and available, and the _outputs_ are known and have an immediate user. You will likely need to work with the user to
ensure that you are automating the solution to a correct problem. Always make sure that all relevant stakeholders agree
on what problem is being solved.

Once inputs and outputs are understood, you can draft some set of **goals** that automation should meet.

### Designing an automated process

An automated solution usually consists of several individual automated processes that may be sequenced or executed in
parallel.

To design a process you have to answer the following questions:
* What are the **inputs** and **outputs**?
* What is the **goal** for this specific process?
* Who would use the **outputs** or side-effects of the process: other automation? some end user?

The **goal** has to be something specific that directly relates to the action performed by the process. For example:
* _Build job_ validates that component can be built
* _Unit test job_ validates that **enabled** unit tests run successfully
* _Integration build_ validates that component can be integrated with other components
* _Integration test_ validates that integrated system behaves as expected

A goal for a single automated process should not be confused with a goal for the whole automated pipeline. Most of the
individual processes do not directly reflect the goals of the automated pipeline, but instead represent the building
blocks that enable the pipeline to achieve its goal. Try not to frame every process as something that contributes to the
goal.

<a name="example_design_flow"></a>
See [Example GOATed design flow](./extra/example_design_flow/) for a hands-on example.

## Putting the GOAT into action

Having a clear purpose for all automation jobs made it easy to communicate with stakeholders and reason about the
unsolicited suggestions. For historic reasons, there was an expectation that CI system could be inherited from the
previous project, however upon closer inspection I had several concerns:
* Every PR build in the legacy system was built using the integration build system (as opposed to standalone builds),
  ultimately running a lighter version of an integrated build, it did the job, but was an overkill for regular PRs
* Every common dependency project had an integrated system build as a PR check in addition to individual build, it was
  common for these build to take more than 3 hours
* Release management utilizing the legacy CI required prolonged code freeze periods to manually align dependencies

Legacy CI evolved to where it was because it was easier to maintain - despite every repository having `Makefile`s or
`CMake` configuration that developers used to build their projects standalone, automation utilized a higher level
`Yocto` build system both for individual and integrated builds. I can see the appeal, but I disagree that automation
behavior for standalone builds should differ from that of developers'. Likewise, I rejected the idea of having
integrated builds as required PR gates for common dependencies, this creates a "moving target problem" when introducing
breaking changes and tightly couples development with integration. It can be a viable approach, but I had strong doubts
that it would work for the current project.

While I was wrapping my mind around the core values of automation, I hacked together some intermediate jobs to satisfy
some urgent requests from opinionated stakeholders. I briefed them about the limitations of these jobs, specifically how
they shape the development process:
* Integrated builds as **required** PR gates slow down the development, as integration of breaking changes would have to
  be done in-sync with the changes themselves, or with external topic-branch synchronization that requires coordinated
  merging of the changes
* Deviating automation behavior from developer actions will result in inconsistent _works-on-my-machine_ PR build
  failures

Sure enough, most of the voiced concerns and risks panned out over the next few months, the _required integrated builds_
became **optional** in about two weeks, and in about 3 months they asked to disable them altogether (they had ~66%
failure rate on PRs with most failures unrelated to the actual change from the PR). Once I finished with the initial
implementation of my automation design, all but one "urgent requests" were phased out completely. The only remaining
legacy feature had to do with some trickery to enable early integration builds to notify stakeholders of failing
components in advance. Thanks to the approach I chose it was relatively easy to translate that work into a separate
automation pipeline specifically for that goal.

All in all, this approach proved to be a success. Even though I had to tweak my design to make it more robust and to
enable the extended set of validation, the core values and the GOAT framework itself did not change at all.

## Eliminating automation SMEs

Ensuring that automation is kept "flat" (most of automated processes are self-contained and there's little crossover
between automated pipelines) and any problem can be pigeonholed into an appropriate pipeline (and then process), ensures
that we don't need "automation subject-matter experts". I have nothing against SMEs in general, but I reject the idea
that everything has to have an SME - automation is not establishing any _new things_ for the project, it merely
automates something that is already known and used. It should not be a complex area that requires some arcane knowledge
of what's-done-where.

The actual SMEs are not the engineers who automated the process, but those who established the original workflow - their
competence should not be hijacked by the person implementing the automated approach. It's very common that person
implementing automation lacks subject-matter knowledge to assess failures of the very process they are automating.

For some reason this was a point of contention in my project, any automated job failure instinctively would be sent for
triage to my team even though we rarely had anything to do with the failure. It took some months to explain that only a
small subset of infrastructure issues can be triaged by my team, the actual observed failures should instead be
redirected to SMEs of the failing components. It seems that management often assumes that responsibility is transferred
together with the automation, even when it's not.

## Healthy ecosystem

This is purely subjective, but I was very happy to see how component SMEs stepped in to improve the automation. Ensuring
that automation does not distance engineers from the delivery provided a new engagement vector for some stakeholders -
they stepped forward with optimizing certain parts of our integration pipelines, they suggested how early integration
activities could be improved, they even developed additional tools to be used by all teams to ensure automation failures
are easier to triage.

Making automation easy to approach and staying open to contributions from developers further minimized the gap between
automation and developer workflows.

## Achievements

The performance of the end system was astounding. There were some challenges on the implementation side to ensure we
took advantage of the build caches, but once we figured it out, other projects reached out to us asking for our lessons
learned and how they could achieve our level of performance:
* **90%** reduction in no-change integrated system build times (from **3 hours** to **15 minutes**)
* **75%** reduction in standalone PR build times (from **2 hours** to **30 minutes**)
* **95%** effort reduction from early integration pipeline requiring no oversight:
  * over 95% of component revision updates being successful
  * 1% of revision updates requiring some developer input to tweak the recipe (e.g. introducing new build-time
    dependencies)
  * 4% of revision updates requiring an extra revision update once the problem is solved within the component repository
  * staying close-to-HEAD required about **1 hour per week** of manual oversight, in contrast to a dedicated engineer
    assigned to manage this for `main` as well as release branches
* **80%** reduction in integrated system release build times (on average reduced from **3 hours** to **40 minutes**)
  enabling daily release candidate iteration, including validation
* **95%** effort reduction from automated release pipeline requiring minimal engineer oversight of **1 hour per RC**, in
  contrast to a dedicated release engineer position
* Engineers relieved by automation were reassigned to tooling and system engineering activities, further improving our
  infrastructure

## Lessons learned

It's a bit challenging to reminisce on this topic, as I don't think I uncovered anything revolutionary or even new. I
ventured into an established field that _was unknown to me_, I understood why and how the CI systems are _the way they
are_. I do not think I arrived at the same destination though. I feel like the system I designed is much more
light-weight and _honest_ with the stakeholders. The documentation for it was also straightforward - it was just stating
what each automated process aims to achieve, what it takes as inputs and what it produces as outputs.

Since the CI system was pretty light-weight (the only complex processes had to do with integration pipelines and their
validation), it was obvious to project management if it satisfied their expectations or not. Whenever there were new
requests to automate something (such as automating tooling migration, or a release process), the communication was very
easy: identify what is different from what's done already, and reframe the request in form of the same old
"input-output-goal" trifecta.

Having the "input-output-goal" framing improved the communications dramatically. A lot is said about "shared
vocabularies" between stakeholders, but I was an outsider in the project.

I asked the opinionated stakeholders about the nature of their requests, I pushed back on the _"just do the same thing
they did for their project"_. This ensured that there was a common understanding of what we're trying to achieve for
_our project_, it allowed us to **avoid mistakes**. It's a common pitfall to jump to known solutions, but you have to
make sure that those solutions actually reflect **your problem**. _The other project_ had their own challenges and
limitations, _our project_ was greenfield and did not inherit either challenges or limitations, but instead we had our
own.

I iterated on the design twice to ensure it accounted for extending the CI system with new "kinds of requests". We also
iterated on the implementation several times, mostly to get rid of the over-abstracted pipelines and hone in on the
actual requests and goals.

I learned a lot in a new area and was very happy with the end result.

### Postscript

I have written separately on the [dark side of automation](../../takes/automation/) and misconceptions surrounding it.
