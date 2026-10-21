package templates

import (
	"fmt"
	"sort"
)

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
	byID map[string]AssetEntry
}

// NewAssetLibrary builds the lookup, rejecting duplicate IDs. IDs are global,
// so two files with the same name in different directories are an error.
func NewAssetLibrary(assets []AssetEntry) (AssetLibrary, error) {
	byID := make(map[string]AssetEntry, len(assets))
	for _, asset := range assets {
		if asset.ID == "" {
			return AssetLibrary{}, fmt.Errorf("asset %s has an empty ID", asset.Path)
		}
		if previous, exists := byID[asset.ID]; exists {
			return AssetLibrary{}, fmt.Errorf(
				"duplicate asset ID %q: %s and %s",
				asset.ID, previous.Path, asset.Path,
			)
		}
		byID[asset.ID] = asset
	}
	return AssetLibrary{byID: byID}, nil
}

// Get returns the asset with the given ID.
func (l AssetLibrary) Get(id string) (AssetEntry, bool) {
	asset, ok := l.byID[id]
	return asset, ok
}

// All returns every asset sorted by ID, for pages that list assets rather than
// render one card.
func (l AssetLibrary) All() []AssetEntry {
	assets := make([]AssetEntry, 0, len(l.byID))
	for _, asset := range l.byID {
		assets = append(assets, asset)
	}
	sort.Slice(assets, func(i, j int) bool { return assets[i].ID < assets[j].ID })
	return assets
}

// Tags returns every distinct tag in sorted order, for the credits filter.
func (l AssetLibrary) Tags() []string {
	seen := make(map[string]bool)
	for _, asset := range l.byID {
		for _, tag := range asset.Tags {
			seen[tag] = true
		}
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
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
