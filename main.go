// Command site builds the Team Filtered website with astwerk.
//
//	go run .            # build into build/
//	go run . -serve     # serve build/ and rebuild on change (drafts on)
//	go run . -drafts    # include devlogs marked as drafts
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/LukasDerBaum42/astwerk/devserver"
	"github.com/LukasDerBaum42/astwerk/ssg"
	"github.com/Team-Filtert/Webside/pages"
	"github.com/Team-Filtert/Webside/templates"
	"github.com/a-h/templ"
)

const outDir = "build"

func main() {
	serve := flag.Bool("serve", false, "serve the build directory and rebuild on change")
	addr := flag.String("addr", devserver.DefaultAddr, "address the dev server listens on")
	drafts := flag.Bool("drafts", false, "include draft devlogs")
	var changed stringList
	flag.Var(&changed, "changed", "changed file for an incremental build (repeatable)")
	flag.Parse()

	if *serve {
		runServer(*addr)
		return
	}

	buildSite(*drafts, changed)
	log.Printf("built the site into %s/", outDir)
}

func runServer(addr string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	templ, templArgs := templCommand()

	// Public images and fonts are copied into the build, so the watcher has to
	// notice them too, not just the default code/content extensions.
	exts := append([]string{}, devserver.DefaultExts...)
	exts = append(exts, ".png", ".jpg", ".jpeg", ".webp", ".gif", ".otf", ".woff", ".woff2")

	ignore := append([]string{}, devserver.DefaultIgnore...)
	ignore = append(ignore, ".gomodcache", ".gocache")

	build := func(changed []string, full bool) error {
		if !full && !dataOnly(changed) {
			full = true
		}

		args := []string{"run", ".", "--drafts"}
		if !full {
			for _, path := range changed {
				args = append(args, "--changed", path)
			}
		}

		return ssg.Steps(
			ssg.Command(templ, templArgs...),
			ssg.Command("go", args...),
		)()
	}

	err := devserver.Run(ctx, devserver.Config{
		Dir:          outDir,
		Addr:         addr,
		Exts:         exts,
		Ignore:       ignore,
		BuildChanged: build,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TemplVersion is the templ release CI pins. Keep them in step.
const TemplVersion = "v0.3.1020"

// templCommand prefers an installed templ binary and falls back to go run,
// matching how CI generates templates.
func templCommand() (string, []string) {
	if _, err := exec.LookPath("templ"); err == nil {
		return "templ", []string{"generate"}
	}
	return "go", []string{"run", "github.com/a-h/templ/cmd/templ@" + TemplVersion, "generate"}
}

func buildSite(includeDrafts bool, changed []string) {
	members := loadMembers()
	assetLibrary := loadAssets(members)
	templates.SetAssetLibrary(assetLibrary)
	log.Printf("loaded %d assets", len(assetLibrary.All()))

	updates := loadUpdates(includeDrafts)
	recent := updates
	if len(recent) > 3 {
		recent = recent[:3]
	}
	templates.SetFooterUpdates(recent)

	updateDates := make(map[string]time.Time, len(updates))
	for _, update := range updates {
		if !update.Published.IsZero() {
			updateDates["updates/"+update.Slug+"/"] = update.Published
		}
	}

	leaders := pick(members, "wolfy", "lukas", "komi", "firearchers", "riko", "admiralensin")
	chaps := loadChapters()

	team := sortedMembers(members)
	groups := teamGroups(team)
	allAssets := assetLibrary.All()
	people := creditPeople(team, allAssets)
	creditGroups, assetsByMember, credited := buildCreditGroups(people, allAssets)
	creditTags := assetLibrary.Tags()

	root := ssg.Node{
		Title: "Neuro-sama: The Good, The Evil And The Filtered",
		// Everything in public/ is copied to the site root, matching how the
		// Astro site served it (img/, assets/, other/, CNAME, robots.txt).
		CopyFrom: "public",
		// Incremental builds rebuild this node when any of these data sources
		// change. The style and scripts directories have their own nodes below.
		DependsOn: []string{"content/members", "content/updates", "public"},
		Page: func(c ssg.Ctx) templ.Component {
			return pages.Home(c, recent, leaders)
		},
		Children: map[string]ssg.Node{
			"about": aboutNode(chaps),
			"team": {
				Title:     "The Team",
				DependsOn: []string{"content/members"},
				Page:      func(c ssg.Ctx) templ.Component { return pages.Team(c, leaders, groups) },
			},
			"contact": {
				Title: "Contact",
				Page:  func(c ssg.Ctx) templ.Component { return pages.Contact(c) },
			},
			"credits": creditsNode(people, credited, creditGroups, creditTags, assetsByMember),
			"updates": updatesNode(updates),
			"style":   {CopyFrom: "style", DependsOn: []string{"style"}},
			"scripts": {CompileFrom: "scripts", DependsOn: []string{"scripts"}},
		},
	}

	// ssg.URLs walks the finished tree (including generated children), so the
	// sitemap cannot drift from the pages Build actually writes.
	root.Files = map[string]ssg.PageFunc{
		".nojekyll": ssg.Static(func() templ.Component { return templ.Raw("") }),
		"404.html": func(c ssg.Ctx) templ.Component {
			// GitHub Pages serves this one file for any missing path, so relative
			// URLs would resolve against that path. The site is deployed at the
			// domain root, so absolute URLs are safe here.
			c.Relative = false
			return pages.NotFound(c)
		},
		"sitemap.xml": ssg.Static(func() templ.Component {
			return ssg.Sitemap(ssg.URLs(root), ssg.SitemapOptions{
				BaseURL: templates.SiteURL,
				LastMod: func(path string) time.Time { return updateDates[path] },
			})
		}),
	}

	if err := ssg.Build(root, ssg.BuildOptions{
		OutDir:       outDir,
		Changed:      changed,
		RelativeURLs: true,
	}); err != nil {
		log.Fatal(err)
	}
}

func aboutNode(chapters []pages.Chp) ssg.Node {
	return ssg.Node{
		Title:     "About the game",
		DependsOn: []string{"content/character.json", "content/chapters"},
		Page:      func(c ssg.Ctx) templ.Component { return pages.About(c, chapters) },
	}
}

func creditsNode(people, credited []templates.Member, groups []pages.CreditGroup, tags []string, assetsByMember map[string][]templates.AssetEntry) ssg.Node {
	return ssg.Node{
		Title:     "Credits",
		DependsOn: []string{"content/members", "public"},
		Page: func(c ssg.Ctx) templ.Component {
			return pages.Credits(c, credited, groups, tags)
		},
		Generate: func() map[string]ssg.Node {
			children := make(map[string]ssg.Node, len(people))
			for _, person := range people {
				person := person
				children[person.Slug] = ssg.Node{
					Title:     "Credits: " + person.Name,
					DependsOn: []string{"content/members", "public"},
					Page: func(c ssg.Ctx) templ.Component {
						return pages.CreditMember(c, person, assetsByMember[person.Slug])
					},
				}
			}
			return children
		},
	}
}

func updatesNode(updates []templates.Update) ssg.Node {
	return ssg.Node{
		Title:     "Updates",
		DependsOn: []string{"content/updates", "public"},
		Page:      func(c ssg.Ctx) templ.Component { return pages.Updates(c, updates) },
		Generate: func() map[string]ssg.Node {
			children := make(map[string]ssg.Node, len(updates))
			for _, u := range updates {
				u := u
				children[u.Slug] = ssg.Node{
					Title:     u.Title,
					DependsOn: []string{"content/updates", "public"},
					Page:      func(c ssg.Ctx) templ.Component { return pages.UpdateDetail(c, u) },
				}
			}
			return children
		},
	}
}
