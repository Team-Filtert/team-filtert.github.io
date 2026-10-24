package main

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/LukasDerBaum42/astwerk/content"
	"github.com/Team-Filtert/Webside/pages"
	"github.com/Team-Filtert/Webside/templates"
)

type memberFrontMatter struct {
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Image       string   `toml:"image"`
	Departments []string `toml:"departments"`
	Language    string   `toml:"language"`
}

type updateFrontMatter struct {
	Title     string   `toml:"title"`
	Summary   string   `toml:"summary"`
	Author    []string `toml:"author"`
	Published any      `toml:"published"`
}

type ChapterFrontMatter struct {
	ID     int    `toml:"chapter"`
	Status string `toml:"status"`
}

// department is one section on the team page. The order here is the order the
// page lists them in.
type department struct {
	ID           string
	Title        string
	ShowLanguage bool
}

var departments = []department{
	{ID: "art", Title: "Meet the Art Department"},
	{ID: "music", Title: "Meet the Audio/Music Department"},
	{ID: "writing", Title: "Meet the Writing Department"},
	{ID: "game-design", Title: "Meet the Game Design Department"},
	{ID: "programming", Title: "Meet the Programming Department"},
	{ID: "translation", Title: "Meet the Translation Department", ShowLanguage: true},
}

// siteLoader is the markdown loader for content that may call templ
// components: {{@asset(id)}} and {{@gallery(["a.png", "b.png"])}}.
var siteLoader = newSiteLoader()

func newSiteLoader() *content.Loader {
	components := content.NewComponents()
	content.Register(components, "asset", templates.Asset)
	content.Register(components, "gallery", templates.Gallery)
	content.Register(components, "castgallery", pages.CastGallery)
	return content.NewLoader(nil).WithComponents(components)
}

func loadMembers() map[string]templates.Member {
	collection, err := content.LoadCollectionWith[memberFrontMatter](siteLoader, "content/members")
	if err != nil {
		log.Fatalf("reading content/members: %v", err)
	}

	members := make(map[string]templates.Member, len(collection.Items))
	for _, item := range collection.Items {
		fm := item.Meta
		if fm.Name == "" {
			log.Fatalf("content/members/%s.md has no name", item.Slug)
		}
		members[item.Slug] = templates.Member{
			Slug:        item.Slug,
			Name:        fm.Name,
			Description: fm.Description,
			Image:       fm.Image,
			Language:    fm.Language,
			Departments: fm.Departments,
			Body:        item.Content,
		}
	}
	return members
}

func loadChapters() []pages.Chp {
	collection, err := content.LoadCollectionWith[ChapterFrontMatter](siteLoader, "content/chapters")
	if err != nil {
		log.Fatalf("reading content/chapters: %v", err)
	}

	chapters := make([]pages.Chp, 0, len(collection.Items))
	for _, item := range collection.Items {
		fm := item.Meta
		if fm.ID == 0 {
			log.Fatalf("content/chapters/%s.md has no chapter number", item.Slug)
		}

		chapters = append(chapters, pages.Chp{
			Slug:   item.Slug,
			ID:     fm.ID,
			Status: fm.Status,
			Body:   item.Content,
		})
	}

	sort.Slice(chapters, func(i, j int) bool { return chapters[i].ID < chapters[j].ID })
	return chapters
}

func loadUpdates(includeDrafts bool) []templates.Update {
	collection, err := content.LoadCollectionWith[updateFrontMatter](siteLoader, "content/updates")
	if err != nil {
		log.Fatalf("reading content/updates: %v", err)
	}

	updates := make([]templates.Update, 0, len(collection.Items))
	for _, item := range collection.Items {
		fm := item.Meta
		if fm.Title == "" {
			log.Fatalf("content/updates/%s.md has no title", item.Slug)
		}

		published, draft, err := publishedDate(fm.Published)
		if err != nil {
			log.Fatalf("content/updates/%s.md: %v", item.Slug, err)
		}
		if draft && !includeDrafts {
			continue
		}

		u := templates.Update{
			Slug:      item.Slug,
			Title:     fm.Title,
			Summary:   fm.Summary,
			Authors:   strings.Join(fm.Author, ", "),
			Body:      item.Content,
			Published: published,
		}
		if draft {
			u.Date = "Draft"
		} else {
			u.Date = published.Format("2 January 2006")
		}
		updates = append(updates, u)
	}

	if len(updates) == 0 {
		log.Fatal("no devlogs found in content/updates")
	}

	sort.SliceStable(updates, func(i, j int) bool {
		if updates[i].Published.Equal(updates[j].Published) {
			return updates[i].Slug < updates[j].Slug
		}
		return updates[i].Published.After(updates[j].Published)
	})
	return updates
}

// publishedDate accepts either a TOML date or the string "draft".
func publishedDate(v any) (time.Time, bool, error) {
	switch value := v.(type) {
	case time.Time:
		return value, false, nil
	case string:
		if value == "draft" {
			return time.Time{}, true, nil
		}
		t, err := time.Parse("2006-01-02", value)
		if err != nil {
			return time.Time{}, false, fmt.Errorf("published: %w", err)
		}
		return t, false, nil
	case nil:
		return time.Time{}, true, nil
	default:
		return time.Time{}, false, fmt.Errorf("published: unsupported value %v", value)
	}
}

// pick returns the named members in the given order, failing loudly on a typo
// rather than rendering a missing card.
func pick(members map[string]templates.Member, slugs ...string) []templates.Member {
	out := make([]templates.Member, 0, len(slugs))
	for _, slug := range slugs {
		m, ok := members[slug]
		if !ok {
			log.Fatalf("no member with slug %q in content/members", slug)
		}
		out = append(out, m)
	}
	return out
}

func sortedMembers(members map[string]templates.Member) []templates.Member {
	slugs := slices.Sorted(maps.Keys(members))

	out := make([]templates.Member, 0, len(slugs))
	for _, slug := range slugs {
		out = append(out, members[slug])
	}
	return out
}

func teamGroups(team []templates.Member) []pages.TeamGroup {
	groups := make([]pages.TeamGroup, 0, len(departments))
	for _, d := range departments {
		var inDepartment []templates.Member
		for _, m := range team {
			for _, id := range m.Departments {
				if id == d.ID {
					inDepartment = append(inDepartment, m)
					break
				}
			}
		}
		groups = append(groups, pages.TeamGroup{
			Title:        d.Title,
			Key:          d.ID,
			Members:      inDepartment,
			ShowLanguage: d.ShowLanguage,
		})
	}
	return groups
}

// creditPeople returns the team plus any external credit referenced by asset
// metadata. External people only have a slug and a display name; their page
// shows the same no-bio fallback as a team member without a bio.
func creditPeople(team []templates.Member, assets []templates.AssetEntry) []templates.Member {
	people := append([]templates.Member(nil), team...)
	seen := make(map[string]bool, len(people))
	for _, member := range people {
		seen[member.Slug] = true
	}

	for _, asset := range assets {
		if asset.Credit == "" || seen[asset.Credit] {
			continue
		}
		seen[asset.Credit] = true
		name := asset.CreditName
		if name == "" {
			name = asset.Credit
		}
		people = append(people, templates.Member{Slug: asset.Credit, Name: name})
	}

	sort.Slice(people, func(i, j int) bool { return people[i].Slug < people[j].Slug })
	return people
}

// buildCreditGroups groups assets by their credit ID. People without
// catalogued assets are left out of the index (but keep their page when they
// are team members), and uncredited assets get a final section of their own.
func buildCreditGroups(people []templates.Member, assets []templates.AssetEntry) ([]pages.CreditGroup, map[string][]templates.AssetEntry, []templates.Member) {
	byMember := make(map[string][]templates.AssetEntry)
	var uncredited []templates.AssetEntry

	for _, asset := range assets {
		if asset.Credit == "" {
			uncredited = append(uncredited, asset)
			continue
		}
		byMember[asset.Credit] = append(byMember[asset.Credit], asset)
	}

	groups := make([]pages.CreditGroup, 0, len(people)+1)
	credited := make([]templates.Member, 0, len(people))
	for _, person := range people {
		if len(byMember[person.Slug]) == 0 {
			continue
		}
		credited = append(credited, person)
		groups = append(groups, pages.CreditGroup{
			Member: person,
			Assets: byMember[person.Slug],
		})
	}
	if len(uncredited) > 0 {
		groups = append(groups, pages.CreditGroup{Title: "Uncredited", Assets: uncredited})
	}

	return groups, byMember, credited
}
