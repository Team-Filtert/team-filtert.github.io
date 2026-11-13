package templates

const (
	// SiteURL is the canonical origin. In-page URLs are relative; only places
	// that must be absolute (canonical/OG tags, the sitemap) use it.
	SiteURL = "https://team-filtered.org"

	SiteTitle = "Neuro-sama: The Good, The Evil And The Filtered"

	SiteDescription = "The official website for Neuro-sama: The Good, The Evil And The Filtered, a Neuro sama fan game made by team filtered."

	DefaultImage = SiteURL + "/img/posters/NGEF_Poster.png"
)

// Meta is the per-page head metadata. Empty fields fall back to the site
// defaults, so the home page can pass the zero value.
type Meta struct {
	Title       string
	Description string
	Image       string // absolute URL
	Type        string // "website" or "article"
}

func (m Meta) TitleTag() string {
	if m.Title == "" || m.Title == SiteTitle {
		return SiteTitle
	}
	return m.Title + " — " + SiteTitle
}

func (m Meta) DescriptionTag() string {
	if m.Description == "" {
		return SiteDescription
	}
	return m.Description
}

func (m Meta) ImageURL() string {
	if m.Image == "" {
		return DefaultImage
	}
	return m.Image
}

func (m Meta) OGType() string {
	if m.Type == "" {
		return "website"
	}
	return m.Type
}
