package templates

import (
	"bytes"
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/LukasDerBaum42/astwerk/ssg"
	templpkg "github.com/a-h/templ"
)

var externalAnchor = regexp.MustCompile(`<a\b[^>]*>`)

// MarkdownBody renders a markdown body with the page's Ctx, so
// content.Page.Render can resolve root-relative URLs and run components.
// External links open in a new tab.
func MarkdownBody(c ssg.Ctx, body func(ssg.Ctx) templpkg.Component) templpkg.Component {
	return templpkg.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		html, err := bodyHTML(ctx, c, body)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, html)
		return err
	})
}

// MemberBio renders a member's bio, with a visible fallback when it is empty.
func MemberBio(c ssg.Ctx, body func(ssg.Ctx) templpkg.Component) templpkg.Component {
	return templpkg.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		html, err := bodyHTML(ctx, c, body)
		if err != nil {
			return err
		}
		if html == "" {
			_, err = io.WriteString(w, `<p class="no-bio">No bio available.</p>`)
			return err
		}
		_, err = io.WriteString(w, html)
		return err
	})
}

func bodyHTML(ctx context.Context, c ssg.Ctx, body func(ssg.Ctx) templpkg.Component) (string, error) {
	if body == nil {
		return "", nil
	}

	var buf bytes.Buffer
	if component := body(c); component != nil {
		if err := component.Render(ctx, &buf); err != nil {
			return "", err
		}
	}

	html := strings.TrimSpace(buf.String())
	return externalAnchor.ReplaceAllStringFunc(html, func(tag string) string {
		if !strings.Contains(tag, `href="http`) || strings.Contains(tag, "target=") {
			return tag
		}
		return strings.TrimSuffix(tag, ">") + ` target="_blank" rel="noopener noreferrer">`
	}), nil
}
