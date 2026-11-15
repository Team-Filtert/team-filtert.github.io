# About page remake — implementation plan

> Technical companion to `About_remake.md`, which holds the page copy.
> Internal only. The chapter teasers below are embargoed until the announcement
> goes live.

## Goals

- Replace the current About page with a clear pitch, a story section, a gameplay
  section, and one card per chapter.
- Keep the work-in-progress status explicit and avoid promises the team cannot
  keep.
- Polish the short About section on the home page at the same time.

## Page structure

1. **Pitch** — what the game is, in three or four sentences.
2. **The story** — the main cast (Neuro, Evil, Nere, Vedal) and the premise.
   Crelly and Ellie move to the chapter cards where they appear, not this
   section.
3. **Gameplay** — Active vs Relaxed combat, normal actions, combo moves,
   ultimates, items and blocks; exploration, quests, shops and inventory as
   planned systems. Source material: Devlog #3.
4. **Chapters** — one card per chapter, loaded from a `content/chapters`
   collection.
5. **Development status** — a short, honest paragraph: Chapter 1 is the demo and
   is in development; nothing has a release date.
6. **Credits** — link to `/credits/`; image credits come from asset metadata.

## Canon and terminology

- **The Filter** — the force Nere embodies and also wields. Public copy can say
  Nere is the embodiment of the Filter and that she is filtering the Neuroverse.
- **The Corruption** — **internal only.** A separate entity that took Nere over
  and explains her behaviour. It is revealed in a secret in-game storyline. Do
  not name it, hint at it, or use "corruption" as a synonym for it on the About
  page or in chapter teasers.
- **Filtered** — the state of someone affected by the Filter. Also the literal
  placeholder text used where copy is not ready yet.
- Nere was released by a drunk Vedal. Keep the comedic tone; treat it as the
  public premise.
- **Eight chapters are planned across seven worlds.** All eight have concepts.
  Chapter 1 is the demo and is in development; the rest are planned.
- Chapters 1 and 8 take place in the same world. **Do not state or imply this
  publicly.**
- Avoid saying "corruption" in the current About copy ("side effect of Nere's
  corruption"); use "the Filter" instead so the term stays unambiguous.

## Chapter data

Put one markdown file per chapter in `content/chapters/`. Load it with
`content.LoadCollection` and sort by `number`.

```toml
+++
number = 2
title = ""                 # optional; the card uses "Chapter N" when empty
status = "planned"         # in-development | planned | released
characters = ["Neuro", "Evil"]   # display names until a characters collection exists
spoiler = "Filtered"       # shown inside <details>; replace when ready
+++

What looks like a normal high school, until something feels off.
```

- The body (everything after the front matter) is the card's overview.
- The card shows the number, status, overview and the confirmed characters.
- An empty body renders the literal text **Filtered**.
- An empty `spoiler` means no `<details>` block is rendered.
- `status = "in-development"` is only for Chapter 1 for now.
- A future `content/characters` collection could replace display names with
  slugs, but it is not required for the first version.
- Add `DependsOn: content/chapters` to the About node in `main.go`.

## Chapter status and teasers (embargoed)

These are the source teasers from the upcoming announcement. Paraphrase for
card copy if needed, but keep the tone and do not add details.

| # | Status | Teaser | Notes |
| --- | --- | --- | --- |
| 1 | in development | No details. Players should go in blind. | Same world as Chapter 8 — internal only. |
| 2 | planned | "What looks like a normal high school, but after a closer look things start to feel off." | |
| 3 | planned | "A Japanese palace where the warriors' primary discipline is… cooking..?" | |
| 4 | planned | "Beyond the stratosphere and into outer space. One small step for Neuro, a big step for AI kind." | |
| 5 | planned | "A circus deep in the woods. While exploring the woods, we find a big mansion — what might await us inside?" | |
| 6 | planned | "A steampunk metropolis in Victorian Ireland. Everywhere you look: cogs, pipes, steam and the mafia…" | Matches Devlog #4's theme change. |
| 7 | planned | "A vast kingdom rife with adventure. Treasures and secrets are just waiting to be found." | |
| 8 | planned | "The game's short but thrilling climax." | Details are censored; do not mention the shared world with Chapter 1. |

Original quotes (internal reference, attributed):

- **[Wolfy]** "As much as I am excited to talk about Chapter 1, we would like to
  refrain from sharing details about it given it is the first to release and we
  want players to go in blind to have a completely fresh experience."
- **[Lukas]** "Chapter 2 takes place in what looks like a normal high school but
  after a closer look things start to feel off."
- **[Wolfy]** "In Chapter 3, the Twins find themselves in a Japanese palace where
  the warriors' primary discipline is… Cooking..?"
- **[Wolfy]** "Chapter 4 takes our AI heroes beyond the stratosphere and into
  outer space! One small step for Neuro, A big step for AI kind."
- **[Lukas]** "Chapter 5 brings us to a circus deep in the woods. While exploring
  the woods, we find a big mansion, what might await us inside?"
- **[Lukas]** "Chapter 6 is set in a steampunk metropolis in Victorian Ireland.
  Everywhere you look are cogs, PIPES, steam and the mafia…"
- **[Wolfy]** "When the Twins land in Chapter 7, they are greeted with a vast
  kingdom rife with adventure. Treasures and secrets are just waiting to be
  found."
- **[Wolfy]** "Chapter 8, the game's short but thrilling climax [censorship] —
  and we really look forward to realizing these visions so that you in the
  audience, devs, and vtubers alike can enjoy the definitive Neuro-sama RPG
  experience."

## Spoiler and placeholder policy

- Before a chapter is released: no plot details. Setting and confirmed
  characters only.
- Where copy is not ready, write the literal text **Filtered** and replace it
  later.
- Use `<details>` for optional spoiler sections. Do not use hover-only reveals;
  they do not work on touch or by keyboard.
- Teaser text is a premise, not a spoiler, and can be shown once the
  announcement is public.

## Copy rules

- **No DLC or postgame promises.** The DLC idea is unlikely and should not be
  mentioned.
- The pitch should mention that the game is open-source and that people are
  welcome to join. Link the Discord and the GitHub org.
- Eight chapters are planned and have concepts, but do not imply dates or order
  of release.
- **Quality:** say we are aiming for a high standard, not that the game will be
  as good as other indie RPGs.
- No release dates, no "coming soon" without a confirmed scope.
- Keep "seven worlds" and "eight planned chapters" consistent between the About
  and home pages.
- Do not reveal the Chapter 1 / Chapter 8 shared world or the Corruption.
- The page copy lives in `About_remake.md`; keep it in sync with these rules.

## Home page About section

- Keep it to two or three sentences: JRPG, Neuro and Evil, seven worlds, the
  Filter, and a short "still in development" note.
- Use the same terminology as the About page.
- Keep the "Learn more..." link.

## Credits

- Replace the hardcoded "Art credits: @TroobsART on twitter" in the Nere card.
  `public/img/characters/meta.json` already has `credit: troobsart` and
  `credit_name: TroobsART` for `nere.jpg`.
- Render the credit from the asset library and link it to
  `/credits/troobsart/`, consistent with the artwork cards.
- Add metadata for the other character images when their credits are known.
- Character cards use their own layout, so read the asset entry rather than
  using `templates.Asset` (which renders a gallery card).

## Implementation checklist

- [ ] `content/chapters/*.md` with the front matter above
- [ ] chapter front matter struct and loader
- [ ] chapter card component
- [ ] `<details>` spoiler component and "Filtered" styling
- [ ] rewrite `pages/about.templ`
- [ ] update the About meta description
- [ ] polish the home page About section
- [ ] read Nere's credit from asset metadata
- [ ] verify credits links and the link checker

## Open items

- Chapter titles: none exist yet, so cards use "Chapter N" until they do.
- Crelly and Ellie: confirm which chapters they appear in before moving them.
- Confirm whether Chapter 5's circus and mansion are one location or two.
- Do not publish any chapter teaser until the announcement is live.
