+++
title = "Devlog #2: Rewrite"
summary = "...yeah we're already rewriting this"
author = ["KTrain"]
published = 2026-06-20
+++


Really soon after the website was launched initially, I was asked to help expand the website in the following ways:

- Adding i18n support (there was only English)
- Making blog posts easier to write (process at the time was just write a HTML page and stick it in, which was not very fun)

There were also some things that griped me from the previous setup:

- The `src/` folder was a mess due to some of Vite's rather annoying limitations.
- Where certain files should be put wasn't obvious at a glance because there was HTML in both `public/` and `src/`.
- There was a bit of unnecessary JavaScript being sent to the client (for the uninitiated, Vite automatically transpiles TypeScript to JavaScript during the build process).
- No uniform way of validating and querying data.
- HTML reuse was nil, so changing something meant remembering to change it across pages.

So, with the task at hand, I determined that a rewrite was already in order.

## The new tech stack

The new version of the website, which you're reading right now, is built with the [Astro](https://astro.build) framework.

The pitch for Astro: it is a static site generator focused around marketing, blogs, personal, and similar sites. The core idea is to make it as easy as possible to ship zero JS to the client, and even if you ship some JS, you shouldn't need to ship a lot. Besides that, the framework also has primitives around working with content and data, and its style of development is easy to pick up and learn. Their docs on Astro components literally says "If you've written HTML, you already know how to write an Astro component!" which is a little bold in my opinion, but probably gets the point across the best way.

I had already had experience with Astro sites - in fact, quite a lot of my projects involved Astro in some way. I like the framework because it perfectly embodies what most sites should be: HTML/CSS first, JS as necessary.

There were also a couple other options I had evaluated, which I'll list in order here:

- Svelte & SvelteKit. I wanted the i18n to not require separate paths and a page load sequence, but the fact I had to ship Svelte to the client made me irk a bit, even if compared to other frameworks it's not quite as heavy. The site shouldn't need to be an SPA, but with this setup, it kinda had to.
- [Content Collections](https://content-collections.dev) + Vite. This felt like a "budget Astro" for most of our purposes. However, using it with HTML pages directly was a bit annoying in my testing, and I didn't want a mess to be left behind with hacking something together.
- Pre-rendered React/Vue/Preact/Svelte (with Vite + [Nitro](https://nitro.build) if necessary). I got some way there, but I thought about it more, and all the extra stuff that Nitro adds can be confusing for team members who aren't as familiar with this style of setup.

Given all that, I felt that Astro was the best choice.

## Migrating

Migrating our pages was fairly simple, given that we had nothing too fancy. The most annoying part was actually migrating some HTML files in `public/` into Markdown files, mostly our members data. This was because the indenting was messed up when I copy-pasted it into a Markdown file. In addition to that, I was trying to replicate the team cards expanding when clicking, but I found that a little hard to do with Astro's content collections and Markdown. Thankfully, Lukas was able to get ahead of me and help fix it and get it over the line (thanks Lukas!).

Our previous setup actually had client-side JavaScript to load some stuff after-the-fact. I think it made sense at that time, but now that we have a framework like Astro, said JS was unnecessary and could be removed. In total, we cut pretty much all of the previous JS code and only left client-side JS for two things:

- Expanding team cards when clicked
- The starry background

The content collections system was really good for us in that each update post, as well as member and department data, had build-time validation on data fields and type-safe accessing as a result. Editor auto-complete was nice too.

The migration is still in progress, and I'll update this post as we migrate more and more. The remaining things left to do, I'm sure we'll get them done soon enough. You'll see it when it happens.
