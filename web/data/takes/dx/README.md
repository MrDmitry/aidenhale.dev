## DX is lyfe

Some rainy day in 2021 I heard the term "developer experience" used in relation to some library holy war. And it dawned
on me that most of my development approaches revolved around that concept! The tools I was developing were all about
the user - the fellow developer (_future me_ most of the time). The way I designed which parts should be
given to the user as parameters, the way I decoupled my tools into modules to be used by other tools - DX was king!

What I didn't grasp at the time was that most of my tools run on the backend and only a small subset of backend
developers use them directly. This changed when we started integrating with another organization. My tools started to
be used widely and the new org had an experienced team of backend developers. These new developers were smart and
**very** opinionated. My tools were a great fit for our integration, so that team embraced them and even started
suggesting improvements. Where it aligned with my vision for the tools, I adopted their proposals. Where it didn't
align, some strongly worded emails were exchanged.

Most of the conflicts were settled with something like "just fork it and maintain a fork for your team - we simply do
not have this use case in our project _or_ let's create a separate tool to cover this difference". And they never did,
I guess in the end it wasn't too important for them.

## DX is a lie

All changed on a rainy spring Monday morning in 2022. I saw an email from one of the opinionated developers from that
new team; the email was sent to both orgs (~400 developers) and had a patch file (!?) attached with a hefty
`-300, +1200` LOC change. But the worst part was not the patch, it was the email.

> This provides a better developer experience

This was the first time DX was used against me. It felt like a betrayal. I thought I was the proud DX enjoyer - surely
they must be wrong. I checked the patch, and of course they were wrong - that's not what **my** tool is supposed to do!

And then I realized it. Well, I realized two things.

1. DX is personal.

    I wouldn't even call it "subjective", DX is how _you_ are used to do things. And then you and your team develop
    some _common_ way of doing things. And then everything that aligns with it becomes "DX" for you.

    DX is the private language you and your community develop for yourself to align on common design and goals. It is
    pretty powerful and it solves a lot of internal friction and conflicts, but step out of your community and DX loses
    its meaning.

2. I wasn't a DX enjoyer.

    I though I developed tools "for DX" but in reality I was just "meeting business goals". _The patch_ hit me with a
    range of end-user workflows that are foreign to how we did things in our CI or deployment. The tools I developed
    first and foremost solved some project needs, and on top of that there was a sprinkle of "this should not be
    painful for others to use".

## The change

I still think about the user, but it's always secondary to the problem I'm solving. It would take a lot of convincing
for me to sacrifice anything for "a better DX". DX changes more frequently than project needs.

I didn't realize it for quite some time, but I designed the tools to play nice with each other, not with the user. I
decoupled modules to be reused to ensure tools behave in a predictable manner and so that I have less code to maintain.

I piped data around, I `tee`'d data and fed it to other tools, I `parallel`'d calls. I did it not because I designed my
tools for it, but because that always met some project goal:
* generate forensics
* keep intermediate files for analysis and recovery
* keep tools simple to maintain
* decouple common code to have less code

If some end-user goals could be solved by combining my tools, I'd rather write a wrapper that pipes data around and
sequences calls in the expected fashion than "extending the tool with some high-level functionality". That's not the
case for prototypes, or first attempts to solve some problem - but when the dust settles and the tool gets rewritten,
there's some modularity, some checkpoints, some interrupts in the flow - things that make it easier to understand the
tool overall.

I changed how I reasoned about my designs, it relieved me of the "will this be a good DX?" conundrum and focused me on
a simpler question of "does this solve the stated problem?". My proposals _didn't change much_, it just became easier
for me to analyze what I was about to prototype and to separate the _actual_ use case from the _nice-to-have_ things.

Did it improve my tools or proposals? Maybe. It definitely reduced scope creep.

## The patch

Did I incorporate _the patch_? Nope. I advised them to fork my tool and they did. As far as I know, they are pretty
happy with their fork.
