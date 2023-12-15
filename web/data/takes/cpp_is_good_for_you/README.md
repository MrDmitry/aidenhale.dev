I'm not going to sell you c++, but I will tell you why c++ was good for me, and why I'm sure it will be good for you.
This is not about the benefits of going low-level, or how c++ skills can translate to other programming languages. This
is about how c++ helps you learn and improve.

### Table of contents
1. [c++ is bad, but it doesn't mean it's _bad_](#c-is-bad-but-it-doesn-t-mean-it-s-bad)
2. [backstory](#backstory)
3. [there is no "c++ way"](#there-is-no-c-way)
4. [there is no "build system"](#there-is-no-build-system)
5. [true sandbox experience](#true-sandbox-experience)
6. [so why is c++ good?](#so-why-is-c-good)

## c++ is bad, but it doesn't mean it's _bad_

> [...] and I think, why you have to be so bad, Zangief? Why can't you be more like good guy? Then I have moment of
clarity... if Zangief is good guy, who will crush man's skull like sparrow's eggs between thighs? And I say, Zangief you
are **bad guy**, but this does not mean you are **bad** guy.
>
> <div class="not-prose"><p style="text-align: right;">Zangief, "Wreck-It Ralph"</p></div>

c++ exposes you to a lot of problems in software development, and more than that, it **delegates** resolving those
problems to the developer. In other languages those are usually solved by built-in tooling, libraries or even some
emergent ecosystem, c++ doesn't really have that and all attempts to do that don't seem to take off.

Exposure to problems makes you a better developer, active work to address the problems makes you a great
developer. So yes, in a sense c++ is good for you because of how bad it is. I'm not talking about "the footguns" or how
low-level it is, I'm talking strictly about how hands-on the c++ developer is with the language.

I hope it's not Stockholm syndrome on my part. For the past few years I explored `rust`, `zig` and `go`, and enjoyed
what I learned. I also got to appreciate c++ more because of it - it's a second nature to me to map the approaches and
concepts I learn in other languages back to c++. Sometimes I find something that doesn't exist in c++ and it leads to an
even more engaging tangent of "how it works", and this fascination and curiosity I attribute to my c++ phase.

## backstory

Becoming a developer in the early 2010s was a unique experience, and I'm not only talking about c++. There were forums
and tutorials, but they were not "bustling with activity". Compsci educational youtube was 90% Indian professors
teaching you by the book, and 10% were some looped royalty-free music over a screengrab of a very slow typist conversing
with you through the magic of `notepad.exe`. Wild stuff.

And so most developers were left to their own devices coming up with creative ways to solve _not-so-complex_ problems by
today's standards. If you were lucky, you had someone to ask. I wasn't very lucky. I found a friend in a CEO of a
startup I was working at the time, and he asked one day _"Can we write backend in c++ instead of php?"_ It could doom
the project since nobody in the startup knew c++, but is it truly a startup without unnecessary risks?

And so a much younger me, with a couple of years of full stack `php+mysql+postres+html+css+js+jquery` behind me,
decided to try some forbidden c++ magic. This was back in 2012, so I equipped myself with
[cplusplus.com](https://web.archive.org/web/20120504003937/http://cplusplus.com/) and a bit later much nicer
[cppreference.com](https://web.archive.org/web/20120722160657/https://en.cppreference.com/w/), and hacked away at all
the segfaults I could catch!

## there is no "c++ way"

The thing I didn't fully recognize at the time, but appreciate now is that there was no "right way" to do it in c++. If
it compiles and does what you expect - it's good! And if it doesn't crash - it's great! And if it doesn't leak resources
after running for a while - you must be pretty good at it already.

It was a great feeling when my multithreaded homebrew sockets worked when I didn't expect them to. They were riddled
with all kinds of race conditions, out-of-order packet processing, out-of-sync data and probably every bug you could
think of. However, that's how I learned about race conditions, memory safety, and how singletons needed to be
initialized before a thread is started or each thread _may_ create its own instance!

Where was I? Oh yes, "the right way". And so you chip away at that leaking and crashing monstrocity of a backend and you
start asking yourself _Is there a better way of doing this?_ And you search the internet and stumble onto some rare gem
from a person who knows what they're doing, like the great
[beej.us](https://web.archive.org/web/20121207020445/https://beej.us/) _(I miss the web design from that era of
internet)_. At first you are overwhelmed with all the new information but soon it starts to make sense. Yet it's not
"the right way", it's the "beej's way", and it's a good way not because it's beej's, but because it works.

And so you appreciate creating things that _just_ work. C++ is already _fast_, so you don't pat yourself on the back
because your program _runs fast_, you pat yourself on the back because it _runs at all_. Eventually you start iterating
faster and get a hang of the compiler and how it works, you learn how to make a better use of the tools you have. And so
you discover the realm of **build tools**.

## there is no "build system"

Love it or hate it, c++ comes with no tools, just a standard and STL implementation. The expectation is that some smart
people will create the tools (because they are good at creating tools) and eventually there will be a whole market of
third-party tools, with all the benefits of competition and blah-blah-blah. Well it kind of happened - you had your
"hero" gcc, "challenger" clang and "the villain" msvc and more.

With compiler sorted out, how do you express your project configuration so you don't compile things by hand? `Makefiles`
of course! Yes-yes, I hear you saying _"But Aiden, people could use Visual Studio 2010 and get all the benefits of
modern IDEs!"_ and sure, but IDEs solve developer problems for us, and that's not what we're after! If we wanted those
problems solved with an IDE, we'd write `java` or `c#`. `Makefiles` were such a pain that when I saw `CMake` the first
time I considered it to be the best thing to ever exist in c++ ecosystem. That says a lot about how bad `Makefiles`
were!

To sum it up, it's not that there are no "build tools", as there are many, there is no _out-of-box_ experience. It's a
choose-your-own-adventure that will expose you to a wide range of problems that exist in software development. And I'd
say modern languages, such as [zig](https://ziglang.org/) are a huge testament to the success of this approach. Would
there be `zig` if c++ had a friendlier syntax and friendlier compilers? Would there be `zig build` if c++ had better
out-of-box DX? There's something to be said about feature bloat of c++, but I'll get to that later.

## true sandbox experience

And so you find yourself with some build tools (that you chose yourself and learned a lot in the process), with some
_requirements_ in mind (e.g. what platforms you want your software to run on, what libraries you intend to use, etc) and
a whole lot of footguns. A lot is written on the topic of footguns, so I'll give a less common take.

Footguns exist in every language, but c++ is **honest** about them.

In fact c++ is so honest about them, they actively introduce features to address them. Fixing footguns is a hard and
long process, there are a lot of things to consider, including the obligatory xkcd's "every change breaks someone's
workflow"

![xkcd strip called "workflow"](https://imgs.xkcd.com/comics/workflow.png)

Jokes aside, I will skim over some of the common misconceptions about c++.

### memory leaks are everywhere

Just pair your `new` with `delete`? Or just use smart pointers and don't worry about clean up. At some point start using
`valgrind`, it's extremely easy to use and incredibly useful.

### pointers are difficult

Just play around with them until it "clicks". It's not _that_ hard, just takes a little bit of getting used to. Software
development is difficult in general, don't be afraid of learning fundamentals.

### c++ is not safe

It's as safe as the code author can prove it to be. There are a lot of tools to assist the author, but with lower-level
languages safety is an acquired skill, not necessarily an out-of-box feature.

### header files!? are we in the 1980s?

There are technical reasons for declaring interfaces outside of implementation, but barring those - it's just a way to
compactly declare your interfaces for the reader (developer, LSP server, compiler). It may feel outdated, but it's truly
not that big of a deal.

### endless feature bloat

Nobody needs to know **all** features that exist. You'll learn the features when you'll need them. And if you didn't
learn them, you either did a good enough job yourself, or didn't need those features at all. When you implement
something yourself, you have a better grasp of what _that thing_ does anyway.

Rarely do you know the problem well enough to jump straight into a working solution, most often you start with a
prototype where you actually learn what you're trying to do. By solving the problem once, you identify what you'd like
to be done differently, and start looking for better ways of doing it. You'll find the answer somewhere, be it in the
STL, `boost`, some random discord server or a separate library that just does what you want.

I agree that feature discoverability is not there, and you need to know what you're looking for. However since you
already solved the problem for your prototype, you should be able to find the relevant things in the STL or elsewhere.

### c++ doesn't have the feature I want to use

If c++ doesn't have some feature and you don't want to implement it yourself - just use some other language. Both `rust`
and `zig` have very handy features for system programming, `go` is great for backend things, `python` is handy for
automation.

### syntax is ugly

Sure, but it's a matter of taste - indeed, reading somebody else's code (even your own code after a couple of weeks,
months, or years) may be unpleasant, but at the end of the day _reading code_ is just a habit. Once you get used to it
you stop noticing the syntax itself.

C++ syntax may feel dated or excessively explicit, but you can style it the way you want it, there are tools for that
too (similar to how you'd use `cargo fmt` or `go fmt`, yet configurable). I do wish for `const` by default though.

### free abstractions are not free

They are free in _some_ sense. Even removing mental overhead often is a good offset for the price you pay.

If you feel like some abstraction is to blame for poor performance - just measure it and refactor accordingly. Often
it's just a skill issue, not necessarily an abstraction cost.

### new standards take years to be supported

Features are usually available as experimental even before the standard is officially adopted, but it does take a while
for _all of the new standard_ to be supported by the mainstream compilers.

However you rarely **need** things from a new standard. Just learn the language, don't get attached to the "latest and
greatest". And if you **need** it, it probably exists in `boost` anyway.

### macros are scary

Just don't use them. Use `constexpr` and `consteval`, they're pretty good and compiler can reason about them.

### dependency management is horrible

Yes. There are some attempts by
[CMake](https://cmake.org/cmake/help/latest/guide/using-dependencies/index.html#downloading-and-building-from-source-with-fetchcontent)
and [conan](https://docs.conan.io/2/introduction.html) to make it better, but it is pretty sad that such a fundamental
area is so underdeveloped.

There are also c++ [modules](https://en.cppreference.com/w/cpp/language/modules) to simplify integration with
dependencies, but even these come with [strings attached](https://youtu.be/_x9K9_q2ZXE?t=2478).

### when is c++2?

[Soon](https://github.com/hsutter/cppfront)

## so why is c++ good?

Instead of focusing on "c++ the language", I will focus on "c++ the development environment".

You get to go as low-level as you'd like. You can specialize in any one area of the language and go deep there. Or you
can stay high-level and just create fast and maybe safe software.

To me it's one of the best system languages there are - love it or hate it but c++ doesn't hold your hand or direct you
to do anything in _some specific way_. All features are opt-in and you get to work directly with all c and c++ libraries
that exist (and being one of the original system languages, there are plenty), no bindings, glue or wrappers necessary.

Even though absense of some "c++ way" leads to big annoyances, like external dependency management and build system
configuration, these are hurdles that you get to overcome **once or twice** and then find and iterate on "your way" of
doing it. There are also many tools helping you bridge these gaps, including `zig build` as your build system, but
that's for a different kind of article.

You get to use debugging and analysis tools that enhance your understanding of software runtime:
* `valgrind` teaches you about memory and profiling
* `gdb` teaches you about runtime and how optimization affects it
* sanitizers and static code analysis tools help you identify errors in your code

You don't **need** to use these when developing, but they are definitely a force multiplier if you do.

I am not saying that c++ should become your primary language, but I am suggesting giving c++ a go. You can start
high-level and then dive deeper wherever you want, at your own pace.

Just remember that there is no "right way" and keep your mind open, it is about the journey after all.
