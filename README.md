# Team Filtered website

Website should be live @ https://team-filtered.org/

Built with [astwerk](https://github.com/LukasDerBaum42/astwerk), a static site
generator for Go + [templ](https://templ.guide/).

## Layout

| Path | What it is |
| --- | --- |
| `main.go` | the `ssg.Node` tree and the build flags |
| `load.go` | `content.LoadCollection`, the department table, and team grouping |
| `assets.go` | the `assets.LoadTree` wrapper that validates sidecar metadata |
| `changes.go` | change tracking for incremental dev builds |
| `sitemap.go` | `sitemap.xml`, `.nojekyll` and the 404 plumbing |
| `templates/` | the shell — layout, nav, footer and cards |
| `pages/` | page bodies |
| `content/updates/*.md` | devlogs, TOML front matter |
| `content/members/*.md` | team bios, TOML front matter |
| `style/` | shared CSS |
| `scripts/` | TypeScript compiled to `build/scripts` |
| `public/` | images and other assets, copied to the site root |

Shared styles stay in `style/`. Page-specific CSS lives in its own style
component next to the page — for example `contactStyles` in
`pages/contact.templ`. `templates.Layout` takes that component as its second
argument and renders it after the footer, at the bottom of `<body>`:

```go
@templates.Layout(c, contactStyles()) {
	...
}
```

Pages without extra styles pass `nil`. `templ` does not scope `<style>` blocks,
so keep selectors unique.

## Dev commands

- `templ generate` — regenerate the `*_templ.go` files after editing a `.templ`
- `go run .` — build into `./build`
- `go run . -dev` — build without cleaning the output directory
- `go run . -serve` — serve `./build` on <http://localhost:8080>, rebuild and reload on change
- `go run . -drafts` — include devlogs marked as drafts
- `go run . -changed PATH` — rebuild only what depends on `PATH` (repeatable)

`-serve` runs `templ generate` and `go run . -dev` for you, so it is the one to
leave running while working. It uses an installed `templ` if available and
falls back to `go run github.com/a-h/templ/cmd/templ@v0.3.1020` otherwise.

Edits under `content/`, `public/`, `style/` and `scripts/` rebuild only the
nodes whose `DependsOn` matches them; anything else — Go, templ, or a deleted
file — triggers a full build.

TypeScript is compiled by `tsc`, so it needs to be installed:

```
npm install -g typescript
```

## Adding content

Bios and devlogs are loaded with `content.LoadCollection`, so the file name is
the slug, TOML front matter is decoded into a struct, and the body is rendered
on demand (`Item.Content`).

A devlog is a markdown file in `content/updates/`. The filename becomes the URL
slug, so `content/updates/Devlog6.md` is served at `/updates/Devlog6/`.

```
+++
title = "Devlog #6: Something happened"
summary = "One line for the index and the home page."
author = ["lukas", "ktrain"]
published = 2026-11-01
+++

Markdown body.
```

`published = "draft"` hides the post unless the build is run with `-drafts`.
Author values are member slugs from `content/members/`.

Devlog bodies may also call templ components. `siteLoader` in `load.go`
registers two:

```
{{@asset("NGEF_Poster.png")}}
{{@gallery(["forrest.png", "tutel_town.png"])}}
```

`asset` renders a single `templates.Asset` card, and `gallery` wraps several of
them in the site's flex gallery. Arguments follow JSON syntax. Asset IDs come
from the `public/` metadata, so credits and alt text stay in one place.

A team member is a markdown file in `content/members/`. The body is the bio and
may contain markdown or raw HTML.

```
+++
name = "KTrain"
description = "Website maintainer"
image = "/img/the_team/KTrain.png"
departments = ["programming", "translation"]
language = "Vietnamese"
+++
```

Departments live in the `departments` table in `load.go`; add a member to a
department in their front matter.

## Assets

Any directory under `public/` can contain a `meta.json` or `meta.toml`. Each
key is a file name in that directory and the asset's global lookup ID; the URL
is generated from the directory, so the asset is never referred to by its full
path. Files without an entry are copied to the build but are not catalogued,
so they cannot be referenced by ID.

```json
{
  "NGEF_Poster.png": {
    "credit": "anzeal",
    "license": "",
    "type": "image",
    "tags": ["poster"],
    "alt": "NGEF key art poster"
  }
}
```

- `credit` is the URL-safe ID used in `/credits/<id>/`. If it matches a member
  slug, the member's name and bio are used. Otherwise it is an external credit:
  `credit_name` sets the display name (falling back to the ID), and that person
  gets a credits page with no bio.
- `type` is `image`, `video`, or `audio`
- `tags` is optional and kept for the credits pages
- `license` is stored for the credits page and is not shown on the card
- `alt` is the image alt text; it falls back to the asset ID

Render an asset with its ID:

```go
@templates.Asset(c, "NGEF_Poster.png")
```

IDs are global: the same file name in two `meta.json` files is a build error.
A missing ID renders a visible error card and logs at build time.

For the credits pages, `templates.Assets().All()` returns every entry sorted by
ID. Each entry carries `Credit` (the member slug), `Tags`, and `License`.

## Credits

- `/credits/` lists everyone with catalogued artwork — team or external — in a
  collapsible list, with tag filters for the artwork
- `/credits/<member-slug>/` shows one member's bio and their catalogued work

Asset cards on these pages use `templates.AssetWithLicense`, which adds the
license when it is set. The credit name on every asset card, the member headings
on the credits index, and the names on the team page all link to the member's
credits page.

## Deployment

GitHub Actions generates the templates, builds the site, and uploads `./build`
to GitHub Pages. The site is served from the domain root, so no base path is
set. All in-page links and assets are relative, so the same build also works
locally and under a subpath.

#### Writing Guide

- Update / Devlogs should always have a Title that starts with "Devlog #x:", where x is the devlog number. Followed by a topic.
