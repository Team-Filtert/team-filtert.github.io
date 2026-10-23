package templates

import (
	"time"

	"github.com/LukasDerBaum42/astwerk/ssg"
	templpkg "github.com/a-h/templ"
)

// Member is a team member loaded from content/members.
type Member struct {
	Slug        string
	Name        string
	Description string
	Image       string
	Language    string
	Departments []string
	Body        func(ssg.Ctx) templpkg.Component
}

// Update is a devlog loaded from content/updates.
type Update struct {
	Slug      string
	Title     string
	Summary   string
	Authors   string
	Date      string
	Body      func(ssg.Ctx) templpkg.Component
	Published time.Time
}

// Character is one entry in content/character.json. The map key is the
// character's name, so loadCharacters fills Name in.
type Character struct {
	Name        string `json:"-"`
	Roll        string `json:"roll"`
	Description string `json:"Description"`
	Pfp         string `json:"pfp"`
	Hiden       bool   `json:"Hiden"`
}

// footerUpdates holds the newest devlogs for the footer. It is installed once
// before Build, like the asset library.
var footerUpdates []Update

// SetFooterUpdates installs the updates shown in the footer. The caller
// decides how many to pass.
func SetFooterUpdates(updates []Update) {
	footerUpdates = updates
}

// leaderLabel keeps the shorter names the home page has always used for two
// of the leaders. The anchor still uses the member's slug.
func leaderLabel(m Member) string {
	switch m.Slug {
	case "wolfy":
		return "Wolfy"
	case "lukas":
		return "Lukas"
	default:
		return m.Name
	}
}

func summaryOr(summary string) string {
	if summary == "" {
		return "No summary available."
	}
	return summary
}

// memberCardID gives the leaders section the stable #member-<slug> anchor and
// everyone else a unique card id, so a member appearing in both the leader and
// department sections does not produce duplicate ids.
func memberCardID(slug, key string, anchor bool) string {
	if anchor {
		return "member-" + slug
	}
	return "card-" + key
}
