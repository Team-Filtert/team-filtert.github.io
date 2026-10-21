package main

import (
	"context"
	"encoding/xml"
	"io"
	"time"

	"github.com/LukasDerBaum42/astwerk/ssg"
	"github.com/Team-Filtert/Webside/templates"
	"github.com/a-h/templ"
)

type sitemapEntry struct {
	Path    string
	LastMod string
}

func sitemap(entries []sitemapEntry) ssg.PageFunc {
	return func(ssg.Ctx) templ.Component {
		return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>`+"\n"); err != nil {
				return err
			}
			if _, err := io.WriteString(w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`+"\n"); err != nil {
				return err
			}
			for _, entry := range entries {
				if _, err := io.WriteString(w, "  <url><loc>"); err != nil {
					return err
				}
				if err := xml.EscapeText(w, []byte(templates.SiteURL+"/"+entry.Path)); err != nil {
					return err
				}
				if _, err := io.WriteString(w, "</loc><lastmod>"+entry.LastMod+"</lastmod></url>\n"); err != nil {
					return err
				}
			}
			_, err := io.WriteString(w, "</urlset>\n")
			return err
		})
	}
}

func pageEntries(urls []string, updates []templates.Update) []sitemapEntry {
	buildDate := time.Now().Format("2006-01-02")

	updateDates := make(map[string]string, len(updates))
	for _, update := range updates {
		if !update.Published.IsZero() {
			updateDates["updates/"+update.Slug+"/"] = update.Published.Format("2006-01-02")
		}
	}

	entries := make([]sitemapEntry, 0, len(urls))
	for _, path := range urls {
		lastMod := buildDate
		if date, ok := updateDates[path]; ok {
			lastMod = date
		}
		entries = append(entries, sitemapEntry{Path: path, LastMod: lastMod})
	}
	return entries
}

// empty writes nothing, for marker files like .nojekyll.
func empty(ssg.Ctx) templ.Component {
	return templ.ComponentFunc(func(context.Context, io.Writer) error { return nil })
}
