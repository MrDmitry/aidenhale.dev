## Automation does not create value

On its own it never creates value for the project, only optimizes the efforts. The value is delivered by the _workflow_
that was originally developed, automation simply **moves the human factor from the manual runtime to its scripted
representation**.

For automation to be valuable, the automated process must:
* be somewhat important
* be used rather frequently
* require significant efforts when done manually

## Automation is not free

Someone has to develop it. If it's not the same person who established the original workflow, it requires more people to
automate a thing, than to do it manually. Automated solution needs to be validated, error checking and recovery options
could also be scripted, status notification mechanisms are usually put in place.

Automation recoups its costs in a funny fashion - on one hand it returns engineer-hours, on the other hand it still
requires engineer-hours to analyze failures and potentially maintain that automation. It's never as simple as "we
automated it and now we never look at it again".

## Automation is not a replacement for an expert

You **need an expert** to develop some workflow. You **need an expert** (could be the same one) to wrap that workflow in
some headless script. The resulting automation does not replace either expert.

Whenever an automated process fails, you still need someone to analyze what went wrong. Depending on the cadence of your
automated process (daily "integration process" vs quarterly "release process") things can go wrong in a variety of ways
and on a very different scale - projects evolve, after all.

You always need someone to triage the newly observed failures. Sometimes it will be caused by the errors in the
automation scripts, sometimes it will be caused by the flaws in the original workflow.

What I'm trying to say is that by automating your processes you're freeing up experts, but you still **need them** to
maintain your automation. Maybe automation allows you to train new experts, but it never replaces them.

## Automation can be a red herring

I've been playing "The Talos Principle 2" and there was a reference to a concept called a [Thought-terminating
cliché](https://en.wikipedia.org/wiki/Thought-terminating_clich%C3%A9) and it rang so many bells in my head, as I've
been a witness to it at the workplace countless times.

To put it simply, there are certain phrases that are simply too agreeable on the surface. So agreeable in fact, that a
listener does not even want to critically evaluate them. "Automation" is one of them.

So many times have I heard "_We can validate it using our automation_" or "_This can be automated_" declared by a
manager on cross-team call. Not once have I heard people push back. Not once was it as simple as "_just run the same
automation we have to validate this_". I started pushing back only when we did it to ourselves a couple of times.

Automation is a very simple concept - a process that does the same thing, usually on some changing inputs. It could be
something as routine as "_running unit tests_" or "_attempting early integration_", or it could be something that does
not usually happen during development - "_preparing a release candidate_".

Whenever a new automated process is integrated into a project, there are always sighs of relief and pats on the backs:
* **engineer** is finally done with it and can switch to something else
* **manager** gets to boast about the newly automated thing
* **stakeholders** get to see yet another indicator on their dashboard

It will work flawlessly, if only for a while, and then some stakeholder notices a light turn red, they ask the manager
why is the light not green, the manager goes to the engineer, and the engineer blames the previous guy (usually
themselves) for missing something obvious. Nobody likes to think about that.

Whenever someone mentions "automation" as a solution, it's very easy to recall that soothing, warm and fuzzy feeling of
"_hey, it works flawlessly_" and disengage from the discussion completely. No harm intended, it is usually an attempt to
keep the discussion from derailing into technical details. Yet, some detailed discussions are necessary.

A stakeholder mentions "automation" in passing, usually expecting someone to double-check them if the proposed existing
automation can be used for the stated problem, or if there is time and bandwidth to develop new automation around the
problem. The listeners disengage because they agree with the _idea of automation in general_. So the stakeholder did not
critically evaluate their own proposal and nobody spoke up - "_I guess everyone agrees_". It's such a powerful word that
it prevented everyone, including the speaker, from critically evaluating the proposal and resulted in everyone nodding
in agreement.

Next time you want to propose "automation", double-check yourself and actually evaluate if existing automation addresses
the problem or if developing a new automated process fits in the stated timeline.

## Bad automation is worse than no automation

Unless you _solved the-problem-at-hand before_, your first solution is probably bad and should be discarded completely,
as you will do a better job on a rewrite.

Same principle is true for automation - your first attempt at automating a new process will most likely be bad. Maybe
even your second and third.

When following some workflow manually you have all your experience at your fingertips. You know the kinds of errors you
typically see _during development_, so it's easy for you to _recognize_ the same errors when following some checklist.
But try to script it and you fall into a rabbit hole of all the potential error states, where they could originate,
which ones you can attempt to retry, when you can ignore some states and when you should fail early.

If you do a bad job with automation, you will just have to maintain it more. And when you're pressed on time to deliver
"something that works", you may be too engaged into fixing it and missing an opportunity for a rewrite. And then it
turns into a [sunk-cost fallacy](https://en.wikipedia.org/wiki/Escalation_of_commitment).

To trivialize, if it takes a comparable amount of time to _maintain an automated process_ as it took to _manually follow
the same workflow_, it's probably bad and you should rewrite it.

## Good automation exists, I think

All that said, if you treat automation not as a buzzword, but as a tool, you can make it work. If you're clear with the
stakeholders about what certain automation achieves, **and what it does not attempt to achieve**, you may be on the
right track!

I have written about my [deep dive with automation](../../work/automation) and the core values I place in automation. I
still had my fair share of errors and rewrites for some processes. I had a year-long stretch of a daily automated
pipeline working flawlessly only for it to spectacularly fail when some unexpected inputs were provided. Nothing is
perfect, and automation should not be assumed to be perfect. It is just another tool that does as good of a job, as you
made it do.
