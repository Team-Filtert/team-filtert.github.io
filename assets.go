package main

import (
	"log"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/LukasDerBaum42/astwerk/assets"
	"github.com/Team-Filtert/Webside/templates"
)

const publicDir = "public"

// assetMeta is one entry in a public/**/meta.json file.
type assetMeta struct {
	Credit     string   `json:"credit" toml:"credit"`
	CreditName string   `json:"credit_name" toml:"credit_name"`
	Alt        string   `json:"alt" toml:"alt"`
	License    string   `json:"license" toml:"license"`
	Type       string   `json:"type" toml:"type"`
	Tags       []string `json:"tags" toml:"tags"`
}

// creditID is the slug used in /credits/<id>/, so it has to be URL-safe.
var creditID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// loadAssets scans public/ for meta.json/meta.toml sidecars and combines every
// catalogued file into one global, ID-indexed lookup. The asset's URL comes
// from the scanner, so it is always relative to the public root.
func loadAssets(members map[string]templates.Member) templates.AssetLibrary {
	entries, err := assets.LoadTree[assetMeta](publicDir, assets.Options[assetMeta]{
		RequireMeta: true,
	})
	if err != nil {
		log.Fatalf("reading %s: %v", publicDir, err)
	}

	catalogued := make([]assets.Entry[templates.AssetEntry], 0, len(entries))
	for _, entry := range entries {
		id := entry.File.Name
		meta := entry.Meta

		switch meta.Type {
		case "image", "video", "audio":
		default:
			log.Fatalf("%s: %q has unsupported type %q", entry.File.Path, id, meta.Type)
		}

		if meta.Credit != "" && !creditID.MatchString(meta.Credit) {
			log.Fatalf("%s: credit %q is not a valid ID", entry.File.Path, meta.Credit)
		}

		// A team member's name and bio come from content/members. Any other
		// credit is an external contributor: credit_name is its display name,
		// and the slug is the fallback.
		creditName := meta.CreditName
		if member, ok := members[meta.Credit]; ok {
			creditName = member.Name
		} else if creditName == "" {
			creditName = meta.Credit
		}

		alt := meta.Alt
		if alt == "" {
			alt = id
		}

		catalogued = append(catalogued, assets.Entry[templates.AssetEntry]{
			File:    entry.File,
			HasMeta: true,
			Meta: templates.AssetEntry{
				ID:         id,
				Alt:        alt,
				Path:       entry.File.Path,
				Type:       meta.Type,
				Credit:     meta.Credit,
				CreditName: creditName,
				License:    meta.License,
				Tags:       normalizeTags(meta.Tags),
			},
		})
	}

	library, err := templates.NewAssetLibrary(catalogued)
	if err != nil {
		log.Fatalf("building asset library: %v", err)
	}
	return library
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			out = append(out, tag)
		}
	}
	sort.Strings(out)
	return slices.Compact(out)
}
