Personal website powered by the **METH** stack.

## The METH stack

<a href="https://daringfireball.net/projects/markdown/" target="_blank">**M**arkdown</a> handles the content and
formatting.

<a href="https://github.com/labstack/echo" target="_blank">**E**cho</a> handles the templating, rendering and general
webserver things.

<a href="https://github.com/tailwindlabs/tailwindcss/" target="_blank">**T**ailwindcss</a> handles the styling.

<a href="https://htmx.org/" target="_blank">**H**tmx</a> handles the frontend.

I could throw another T in there, as I am using [Toml](https://toml.io/en/) for article metadata in the absence of a
database, but I feel like the acronym is already good enough.

## Return to monke

I wanted to create a blog. I worked in webdev in the early 2010s but somehow skipped most of the rapid evolution that
happened with react, django, hugo and all the others. Over the past few years I played around with different frameworks
and tools (angularjs, react, vue, django, nextjs, hugo) and found most of them either overwhelming, or an overkill for a
personal blog or a hobby project. So I figured I could combine developing my personal website with learning some go.

I already had a collection of Markdown documents with notes I drafted over the years, some in better condition than
others, and I wanted to upload them to GitHub for sharing them with friends and colleagues. But I just could not stop
myself from using this opportunity to also learn some go along the way.

I decided to reject the mainstream stacks and return to monke - use filesystem as my database, cache the articles in go
runtime and render Markdown to HTML on the server side.

Tailwind was an easy pick, as it felt like I'm back in the <a
href="https://web.archive.org/web/20131202065213/http://getbootstrap.com/examples/theme/" target="_blank">good ol'
bootstrap days</a>.

The only remaining part was to figure out how interactive (if at all) the blog should be. I could go the static page
route (and coincidentally learn how to prepare static websites) or the dynamic route, with filtering, infinite scrolls
and all that jazz. How could I pass the opportunity to use `htmx`?! Of course it has to be dynamic!

Even though I'm not using most of htmx features, the whole approach felt very refreshing. It combines server rendering
approach with reactive UI and seamless history management without you writing a single client line in js - it does it
all for you, as long as you annotate it in html. It took a bit to get used to it, but once it clicked it felt awesome.

## License

This website is dual-licensed under MIT and CC-BY-4.0. All content is licensed under MIT, non-source code materials are
additionally licensed under CC-BY-4.0.

`SPDX-License-Identifier: MIT AND CC-BY-4.0`

