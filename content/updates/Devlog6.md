+++
title = "Devlog #6: Yet another rewrite..."
summary = "I did another rewrite of the website. We made progress on the actuall game."
author = ["LukasDerBaum"]
published = "draft"
+++

## The website rewrite

Yes, I know. Another rewrite. This time it's the website (again), not the game.

The site has been running on [Astro](https://astro.build) since [Devlog #2](/updates/Devlog2/), and that was a big step up from the original setup. But it still meant keeping a whole Node toolchain around to build a mostly static site, and every small change went through pnpm and a pile of dependencies.

So the site is now built with [astwerk](https://github.com/LukasDerBaum42/astwerk), my own static site generator written in Go. Pages are Go and templ components, devlogs and bios are still Markdown with TOML front matter, and `go run .` builds the whole thing. Astro and pnpm are gone. The few interactive bits are still TypeScript, because I'm not rewriting the starfield in Go again (yes, I did do that for my own website, and it's a bad idea).

The content didn't change, but a few things got nicer along the way:

- Artwork is catalogued once in `public/**/meta.json` — credit, license, type, tags and alt text — instead of being hardcoded into each page.
- That catalogue builds the credits pages. Everyone with catalogued artwork gets a page, and the overview can be filtered by tag.
- Devlogs can drop in artwork with a component call like `{{@gallery(["forrest.png", "tutel_town.png"])}}`. The image path, alt text and credit all come from the metadata, so there is nothing to copy by hand.
- The dev server only rebuilds the pages a change actually affects, which makes writing these posts a lot faster.

Deployment is still GitHub Pages, so nothing changes on your end. It just takes one fewer toolchain to get there.

### Other changes to the site

Together with the rewrite we went allong and made a few more chages to the site as well as adding new stuff.

One of though changes was the redesing of the [About](/about/) page. The old one didn't really have much to show so the new one has much more content
and the content is much better structured.

Another addition is the new [Credits](/credits/) page with contains the credits for all the stuff on the site and Later the assets of the game.


## The Game progress

We have made huge progress with the combat system, after porting most of the old system to the new codebase.
Yes in the end the rewrite the the game wasn't the best idea but we're getting close to something where we can actually so something off.

## Astwerk

<details>
<summary>Let me talk about astwerk a bit.</summary>

Astwerk is a static site generator for Go and [templ](https://templ.guide/). I built it because I really don't like JavaScript and the Node.js ecosystem, but I want to do web development in a language I actually enjoy, with a build that is just `go run .` and a framework that gives you total freedom and control to do things the way you want. It's still pre-1.0, and this site is one of the first to use it. The others are the [Astwerk website](https://lukasderbaum42.github.io/astwerk-website/) and my personal website, so most of it gets used and abused here as a form of stress test.

The basic model is a tree of nodes written in Go. A node can be a page, a directory, a folder copied from `public/`, or a set of pages generated from a collection. Children hang off their parent, and the shared chrome — layout, nav, footer — is just templ components. There is no file-system routing convention to learn: the tree in `main.go` is the site.

The rest is the stuff a site needs, without me gluing it together by hand:

- Markdown content with typed front matter
- Asset metadata scanned from sidecar files, with one global ID per asset
- Relative URLs, so the same build works at the domain root or under a subpath
- A sitemap generated from the same node tree
- A dev server that watches files and rebuilds only the nodes affected by a change

Markdown can also call templ components now, which is how the galleries in earlier devlogs are rendered. That was the last piece that made writing these posts in Markdown feel as capable as writing them in HTML.

It's not trying to be Astro or Hugo. It's a small, opinionated generator shaped around my personal preferences. The source is on [GitHub](https://github.com/LukasDerBaum42/astwerk) if you want to poke at it.

There are also some other cool features planned, some of which will live in their own separate addon repos:

- A theming system with its own data language
- A system that lets you write Go for the front end
- A preprocessor similar to Svelte, for reactive components written in Go

</details>