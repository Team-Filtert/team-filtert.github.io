// Command site builds the Team Filtered website with astwerk.
//
//	go run .            # build into build/
//	go run . -dev       # keep the existing output and overwrite in place
//	go run . -serve     # serve build/ on localhost:8080 and rebuild on change
//	go run . -drafts    # include devlogs marked as drafts
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/LukasDerBaum42/astwerk/devserver"
	"github.com/LukasDerBaum42/astwerk/ssg"
	"github.com/Team-Filtert/Webside/pages"
	"github.com/Team-Filtert/Webside/templates"
	"github.com/a-h/templ"
)

func main() {
	dev := flag.Bool("dev", false, "keep the existing output and overwrite in place")
	serve := flag.Bool("serve", false, "serve the build directory and rebuild on change")
	addr := flag.String("addr", devserver.DefaultAddr, "address the dev server listens on")
	out := flag.String("out", "build", "output directory")
	drafts := flag.Bool("drafts", false, "include draft devlogs")
	var changed stringList
	flag.Var(&changed, "changed", "changed file for an incremental build (repeatable)")
	flag.Parse()

	if *serve {
		runServer(*out, *addr, *drafts)
		return
	}

	if err := buildSite(*out, *dev, *drafts, changed); err != nil {
		log.Fatal(err)
	}
	log.Printf("built the site into %s/", *out)
}

func runServer(outDir, addr string, drafts bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	templ, templArgs := templCommand()
	tracker := newChangeTracker()

	// Public images and fonts are copied into the build, so the watcher has to
	// notice them too, not just the default code/content extensions.
	exts := append([]string{}, devserver.DefaultExts...)
	exts = append(exts, ".png", ".jpg", ".jpeg", ".webp", ".gif", ".otf", ".woff", ".woff2")

	ignore := append([]string{}, devserver.DefaultIgnore...)
	ignore = append(ignore, ".gomodcache", ".gocache")

	build := func() error {
		changed, removed, initial, err := tracker.diff()
		if err != nil {
			return err
		}
		if !initial && !removed && len(changed) == 0 {
			return nil
		}

		full := initial || removed || !tracker.dataOnly(changed)
		args := []string{"run", ".", "--out", outDir, "--dev"}
		if drafts {
			args = append(args, "--drafts")
		}
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
		Dir:    outDir,
		Addr:   addr,
		Exts:   exts,
		Ignore: ignore,
		Build:  build,
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

func buildSite(outDir string, dev, includeDrafts bool, changed []string) error {
	members, err := loadMembers()
	if err != nil {
		return err
	}
	assetLibrary, err := loadAssets(members)
	if err != nil {
		return err
	}
	templates.SetAssetLibrary(assetLibrary)
	log.Printf("loaded %d assets", len(assetLibrary.All()))

	updates, err := loadUpdates(members, includeDrafts)
	if err != nil {
		return err
	}
	if len(updates) == 0 {
		return errors.New("no devlogs found in content/updates")
	}

	leaders, err := pick(members, "wolfy", "lukas", "komi", "firearchers", "riko", "admiralensin")
	if err != nil {
		return err
	}
	heads, err := pick(members, "wolfy", "lukas", "komi", "firearchers", "riko", "admiralensin")
	if err != nil {
		return err
	}
	groups := teamGroups(members)

	team := sortedMembers(members)
	allAssets := assetLibrary.All()
	people := creditPeople(team, allAssets)
	creditGroups, assetsByMember, credited := buildCreditGroups(people, allAssets)
	creditTags := assetLibrary.Tags()

	latest := updates[0]

	root := ssg.Node{
		Title: "Neuro-sama: The Good, The Evil And The Filtered",
		// Everything in public/ is copied to the site root, matching how the
		// Astro site served it (img/, assets/, other/, CNAME, robots.txt).
		CopyFrom: "public",
		// Incremental builds rebuild this node when any of these data sources
		// change. The style and scripts directories have their own nodes below.
		DependsOn: []string{"content/members", "content/updates", "public"},
		Page: func(c ssg.Ctx) templ.Component {
			return pages.Home(c, latest, leaders)
		},
		Children: map[string]ssg.Node{
			"about": {
				Title: "About the game",
				Page:  func(c ssg.Ctx) templ.Component { return pages.About(c) },
			},
			"team": {
				Title:     "The Team",
				DependsOn: []string{"content/members"},
				Page:      func(c ssg.Ctx) templ.Component { return pages.Team(c, heads, groups) },
			},
			"contact": {
				Title: "Contact",
				Page:  func(c ssg.Ctx) templ.Component { return pages.Contact(c) },
			},
			"credits": {
				Title:     "Credits",
				DependsOn: []string{"content/members", "public"},
				Page: func(c ssg.Ctx) templ.Component {
					return pages.Credits(c, credited, creditGroups, creditTags)
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
			},
			"updates": {
				Title:     "Updates",
				DependsOn: []string{"content/updates", "content/members"},
				Page:      func(c ssg.Ctx) templ.Component { return pages.Updates(c, updates) },
				Generate: func() map[string]ssg.Node {
					children := make(map[string]ssg.Node, len(updates))
					for _, u := range updates {
						u := u
						children[u.Slug] = ssg.Node{
							Title:     u.Title,
							DependsOn: []string{"content/updates", "content/members"},
							Page:      func(c ssg.Ctx) templ.Component { return pages.UpdateDetail(c, u) },
						}
					}
					return children
				},
			},
			"style":   {CopyFrom: "style", DependsOn: []string{"style"}},
			"scripts": {CompileFrom: "scripts", DependsOn: []string{"scripts"}},
		},
	}

	// ssg.URLs walks the finished tree (including generated children), so the
	// sitemap cannot drift from the pages Build actually writes.
	root.Files = map[string]ssg.PageFunc{
		".nojekyll": empty,
		"404.html": func(c ssg.Ctx) templ.Component {
			// GitHub Pages serves this one file for any missing path, so relative
			// URLs would resolve against that path. The site is deployed at the
			// domain root, so absolute URLs are safe here.
			c.Relative = false
			return pages.NotFound(c)
		},
		"sitemap.xml": sitemap(pageEntries(ssg.URLs(root), updates)),
	}

	return ssg.Build(root, ssg.BuildOptions{
		OutDir:       outDir,
		Dev:          dev,
		Changed:      changed,
		RelativeURLs: true,
	})
}
