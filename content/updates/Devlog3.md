+++
title = "Devlog #3: State of the demo"
summary = "An update on our progress towards releasing the demo (chapter 1)"
author = ["lukas", "firearchers", "wolfy", "gab"]
published = 2026-07-08
+++


## General
### Gameplay
Most of the gameplay mechanics like
- NPCs
- Combat
- Items

are implemented, but we are currently working on refactoring of the codebase so progress with actually putting things in the game has been slow.

(yes I know a lot of refactoring/rewriting going on here at team Filteres)

#### Lets take about some of our there plans:

- We're planning on having Shops in which the player can buy or sell items.

- We have plans for a bunch of interesting side quests to complete. Those side quests will have give special rewards and unlock new content. (Some may even impact the game world and story in some way)

- There will be an inventory system in which the player can store, equip and use the items they collect. Those items can give new abilities or improve existing ones and a character's stats.

### Combat
The combat system in mostly implemented, with mostly attacks and a few special abilities missing.

But can give a rundown of the combat system is like:

The Combat it as you may know turn-based with inspiration from modern and traditional RPGs.

#### Combat Modes
There are two main combat modes the player can choose from depending on their preference:

- Active Mode: In Active Mode the player can directly influence combat via quick time events. This is similar to "modern" RPGs like Expedition 33, Mario and Luigi or Paper Mario.
This allows a skilled player to deal more damage than a relaxed player but also with the risk of missing.
- Relaxed Mode: In Relaxed Mode the player can passively react to combat, without direct control over the outcome of actions. This is more akin to "traditional" turn-based RPGs. Like Baldur's Gate, Pokemon or Persona.
This mode allows a relaxed player to focus more on strategy with more predictability.

The player can choose in both modes, what action they party member should take and what the target of that action is.

*Note:* The Relaxed Mode is also there to allow Neuro and Evil to play the game, as we plan on adding Neuro integration when the full game is out. 

#### Combat Actions
There are several types of combat actions available to the player:

- Normal Actions: Normal Actions can be, basic attacks or abilities that buff or debuff one or more targets. 
- Combo Moves: Combo Moves are actions performed by 2 or more party members together. (this is largely inspired by the Mario and Luigi series)
- Ultimate Moves: Ultimate Moves are powerful moves that charge up over time and require a full charge to be used. (Wolfy wanted them)
- Items: You mechanics are consumable items that can have several effects. (e.g. healing, buffing, debuffing, or even special attacks)
- Blocks: Blocks are automatically triggered when attacked, and can reduce the amount of damage taken and maybe even reflect some of it back to the attacker.


### Art
Our main bottle neck right now is art.
Our codebase is missing most of  the art assets needed for the game.

At the moment, we're using placeholder assets to fill in some of the gaps.

Art that is missing:
- Character art
- Environment / Tilesets
- Combat sprites
- UI assets
- Item sprites

If you're an artist or a designer, you're welcome to join us and help us with the art assets.

A new organizational system has been made to keep art flowing, as many things in the game are missing visuals to accompany them.


### About the rewrite
Hey, this is Gab, one of the programmers here at Team Filtered.

I wanna start off by saying that I joined the project six-ish months ago knowing literally nothing about game development. Some would probably call me a crazy person for taking on something of this scope as my first game, and honestly, they might have a point. The swarm's collective dedication is just that contagious, what can I say.

Anyways, we're currently in the middle of an architecture rework, which requires a lot of design and thinking before you actually write any code (found out the hard way). A lot of trial and error has been involved, as well as an unholy amount of research. I've personally spent dozens of hours reading forums, documentation, Discord servers, and even books (crazy, I know). I admit it's a bit frustrating at times, but I feel that I've learned a lot.

Lately, my main focus has been on designing and testing the core systems of the game, more specifically, how we handle level switching and transition effects. Those two things as a pair were particularly painful to get right because of my flaky understanding of how coroutines work in Godot and because I kept getting my signals all tangled up. I think I've landed on something solid for the time being, so I'm happy about that.

To be honest, this stage of development is a bit boring, so there aren't any interesting screenshots or game mechanics to show yet. But I believe that getting this right will allow us to start picking up the pace and actually focus on the fun stuff our (very creative) writing department has planned.

That's all from me. Thanks for reading.

Want to help us build this thing?

We're always looking for new contributors. If you've got skills relevant to writing, game design, art, programming, or really anything you think might be useful (and tolerance for a bit of scuff), you can find us over on our [Discord](https://discord.gg/nAXzxMY9gT) and get involved.

## Chapter 1

Development has been slow as of late, but some important progress towards the Chapter 1 demo has been made.

### The Story
The main story for Chapter 1 has been outlined and is largely complete. The remaining work focuses on refining smaller story details like codex entries and writing dialogue.

### NPCs
The core concepts and roles for the Chapter 1 NPCs have been planned. The next step is implementing them in-game and completing their dialogue.

### Quests
The quests for Chapter 1 have been designed conceptually. Implementation and playtesting are the primary tasks remaining.

### Items
All Chapter 1 item concepts have been designed and are ready for implementation.

### Game Design
Most of the game design for the demo is already in place. The map layout is largely complete, with only the final environment textures still to be adde
