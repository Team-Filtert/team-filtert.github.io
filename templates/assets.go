package templates

import "github.com/LukasDerBaum42/astwerk/assets"

// AssetEntry is one catalogued file described by a public/ meta.json.
type AssetEntry struct {
	ID         string
	Alt        string
	Path       string // site-relative, e.g. "img/posters/NGEF_Poster.png"
	Type       string // image, video, audio
	Credit     string // member slug
	CreditName string // display name, resolved when the index is loaded
	License    string
	Tags       []string
}

// AssetLibrary is the global, ID-indexed asset lookup.
type AssetLibrary struct {
	index *assets.Index[AssetEntry]
}

// NewAssetLibrary builds the lookup, rejecting duplicate IDs. IDs are global,
// so two files with the same name in different directories are an error.
func NewAssetLibrary(entries []assets.Entry[AssetEntry]) (AssetLibrary, error) {
	index, err := assets.NewIndex(entries, assets.IndexOptions[AssetEntry]{
		Key:  func(entry assets.Entry[AssetEntry]) string { return entry.Meta.ID },
		Tags: func(entry assets.Entry[AssetEntry]) []string { return entry.Meta.Tags },
	})
	if err != nil {
		return AssetLibrary{}, err
	}
	return AssetLibrary{index: index}, nil
}

// Get returns the asset with the given ID.
func (l AssetLibrary) Get(id string) (AssetEntry, bool) {
	if l.index == nil {
		return AssetEntry{}, false
	}
	entry, ok := l.index.Get(id)
	return entry.Meta, ok
}

// All returns every asset sorted by ID, for pages that list assets rather than
// render one card.
func (l AssetLibrary) All() []AssetEntry {
	if l.index == nil {
		return nil
	}
	entries := l.index.All()
	out := make([]AssetEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.Meta)
	}
	return out
}

// Tags returns every distinct tag in sorted order, for the credits filter.
func (l AssetLibrary) Tags() []string {
	if l.index == nil {
		return nil
	}
	return l.index.Tags()
}

var assetLibrary AssetLibrary

// SetAssetLibrary installs the build-time lookup. Call it once before Build;
// the map is read-only afterwards.
func SetAssetLibrary(library AssetLibrary) {
	assetLibrary = library
}

// Assets returns the installed library for pages that list assets.
func Assets() AssetLibrary {
	return assetLibrary
}
