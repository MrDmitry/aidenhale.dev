I grew to love linters back in my c++ days. How could you not - they spot errors before the compiler! Then you add them
to CI because linters are amazing, and even expand with more static code analysis behemoths (i.e. SonarQube) to detect
even more things! And you also promote formatters, and configure them the way **you** like to win all the aesthetics
debates!

And then something clicked in me. Linters evolve and update and just start getting in the way. A month passes and your
project adopts another set of policies that get linted and the code you've been working for a while now starts to report
warnings (a lot of them false positives) and CI starts blocking your long-running PRs. Over time your code is littered
with cryptic ignore comments with some policy numbers. And then the linter craze passes, and you just freeze your set of
policies and never touch them again.

A new hire in my recent project suggested we adopt linters for our python tools (we had none) to improve them and my
reflex was to immediately push back, but his tenacity prompted some introspection in me. I decided to figure out what
was my relationship with linters at the end of the day.

## the elitist phase

Looking back at my _#linter4lyfe_ days, I see why I loved them so much - they saved me time by continuously checking my
code. I didn't need for a compiler to tell me if something was wrong - the editor just let me know immediately.

I was also one of the early adopters of linters in my project, so I was excited to share it with my colleagues: _"Look!
No need to hit compile! It just tells you what's wrong or questionable!"_ So naturally I felt like **everyone** needed
to have linters enabled, otherwise you're just slowing yourself down.

What I didn't consider at the time, was the cost of adopting linters.

I spent a couple of days settling on a good set of policies that I could enable for out project, I even worked overtime
to address all of the errors on our main branch, and I was happy to present a modest `-2k, +2k` LOC change to enable it
on CI. Our architect was onboard and we enabled it right away.

You probably guessed, but my `-2k, +2k` change created merge conflicts for everyone. Not everyone was a fan, especially
developers who didn't know how to enable linters in the editor - they had to run linters separately, like they did with
a compiler. You may say _"Skill issue"_, but there's no denying that linter adoption was disruptive to our project.

At the time I felt righteous - I was almost single-handedly improving product quality by using broader set of tools.

## self-confrontation

I was reviewing some PR one day, and was a bit confused by a code snippet. I asked the developer for clarifications, as
it looked like an overcomplicated way of doing a rather simple thing.

_"Linter told me to rewrite something, and this passes the linter"_ was the response.
