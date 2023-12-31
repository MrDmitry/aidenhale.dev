## Background

Consider a system that consists of 3 components:
* `libcommon`
* `appProducer`
* `appConsumer`

We will start with 2 pipelines supporting this system:
1. _Component validation_ ran at component level
2. _System integration_ ran at system level

It may be tempting to immediately intertwine the two pipelines and say that _"component validation enables system
integration"_, but we should start with a simple system and only make it more complex if it provides value.

## Component validation

Pipeline description:
* **Inputs**: version controlled sources
* **Outputs**: build status notifications, generated artifacts
* **Goal**: enable standalone component development and prove that `main` branch is stable (can be built and passes
  desired tests)
* **Consumer**: engineers, potentially _system integration_ automation pipeline

This pipeline consists of the following automated processes:
* **PR builds**
    * **Inputs**: _compare_ branch, _base_ branch
    * **Outputs**: build status, build logs, PR vote, generated artifacts
    * **Goal**: ensure proposed change is buildable and can produce artifacts
    * **Consumer**: PR vote, PR author, _PR tests_ process
* **PR tests**
    * **Inputs**: build artifacts, test inputs
    * **Outputs**: test status, test logs
    * **Goal**: ensure proposed change doesn't break _enabled_ tests (scoped for PR validation)
    * **Consumer**: PR vote, PR author, code quality tracker, test tracking software
* **Scheduled build** (e.g. _nightly_ or at some custom cadence)
    * **Inputs**: `main` branch
    * **Outputs**: build status, build logs, generated artifacts
    * **Goal**: ensure `main` branch is stable
    * **Consumer**: repository owner and maintainers, _Scheduled tests_ process
* **Scheduled tests**
    * **Inputs**: build artifacts, test inputs
    * **Outputs**: test status, test logs
    * **Goal**: ensure component behaves as expected (broader test scope)
    * **Consumer**: repository owner and maintainers, code quality tracker, test tracking software

Note that these processes do not contribute anything to the _system integration_ pipeline, they are focused solely on
standalone component quality. _Scheduled_ processes indirectly complement _system integration_ pipeline by performing
standalone validation, as opposed to end-to-end/system validation, providing additional avenue of investigations in case
failures are identified during system validation.

## System integration

Integration model depends on project management, but ultimately it boils down to ensuring some `latest` and `stable`
combination of components are used together. For simplicity, we'll assume that we have some manifest where we specify
component revisions to be included in the integrated system.

Pipeline description:
* **Inputs**: integration manifest, generated artifacts and/or sources
* **Outputs**: build status notification, generated artifacts
* **Goal**: validate the integration manifest, prove feature delivery
* **Consumer**: product owner and key stakeholders

This pipeline is similar to the component validation, except for the fact that it revolves around integration manifest
and configurations, but otherwise the steps and goals are the same.

We know what actions are to be taken in order to move "system integration" forward - update components' revisions
ensuring that a stable configuration can be achieved. If we encounter a failed system build or tests, we need to bisect
the revision(s) of failed component(s) until a stable configuration is achieved. If no stable configuration can be
achieved - we have an escalation path to notify the responsible stakeholders.

Knowing all the steps of system integration maintenance, we know how to automate it - perform exactly the same steps in
an automated manner. The cadence is not as important - it may be a daily, weekly or even a rolling job that just runs
all the time. What's important is that we deliver exactly what we expect - a stable configuration of integrated
components according to the enabled tests.

## Summary

At the end of this exercise we have a CI system that contributes significant value to the overall project delivery:
* _Component validation_ pipeline ensures that any individual component has sufficient (component-controlled) level of
  validation
* _System integration_ pipeline ensures that there's a `stable` configuration of `latest` (or as close as possible)
  components
* _System integration_ also is self-driven, where process of cherry-picking a `stable` configuration is also automated

We didn't _start_ with a goal of automating integration, but we've arrived there naturally. Each step of the way was
focused on delivering some immediate value - first enabling standalone component development, then ensuring that we can
have a stable configuration. This maps nicely to what would be asked of us in a real project - creating PR jobs to
ensure repositories are healthy, creating system builds, etc.

It is possible to achieve the same processes without this analysis, simply following the common sense of _there has to
be a PR build_ and _there has to be a PR test_, and so on. However focusing on the _solution_ doesn't improve your
understanding of the _problem_. It's very easy to fall into a trap of "it's automated, I don't need to worry about it",
however you will always be behind _whenever_ automation goes wrong. Focusing on the goals and consumers ensures that
there's a known value that automation delivers. When there's some automation failure, we know exactly how we arrive to
the failure point.

## Failure analysis

Consider that a `stable` configuration was not found, and key stakeholders received a rare failure notification. We know
exactly what stage of what pipeline in our CI system feeds into this failure. We know how individual components are
validated, we know how integration attempts are performed, it should be easy for us to discover what stakeholder
expectations were not met.

Let's assume that the failure is caused by inconsistent version of `libcommon` used between `appProducer` and
`appConsumer`. We would know if dependencies are considered by integration maintenance, and therefore could immediately
say how automation could be improved. Addressing the failure also becomes relatively straightforward - we could utilize
dependency information when looking for a `stable` configuration, or we could find any other way of covering that missed
_goal_.

Alternatively, if we did not establish _goals_ for our automated processes, it will not be obvious where the failure
originated. We may be tempted to trigger early system integration on `libcommon` PR level to _identify when we're
breaking downstream components_ (surely this is not taken from a real life project), or we may be tempted to partially
freeze some components (e.g. `libcommon`) from their revision being updated automatically (also surely not from a real
project). Knowing what you **are** delivering with automation is very valuable to make decisions on how to address
failures.

