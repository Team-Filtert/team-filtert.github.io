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

// MemberBio renders a member's markdown body with the page's Ctx, so
// content.Page.Render can resolve root-relative URLs and run components.
// External links open in a new tab, and an empty body gets a visible fallback.
func MemberBio(c ssg.Ctx, body func(ssg.Ctx) templpkg.Component) templpkg.Component {
	return templpkg.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		var buf bytes.Buffer
		if body != nil {
			if component := body(c); component != nil {
				if err := component.Render(ctx, &buf); err != nil {
					return err
				}
			}
		}

		html := strings.TrimSpace(buf.String())
		if html == "" {
			_, err := io.WriteString(w, `<p class="no-bio">No bio available.</p>`)
			return err
		}

		html = externalAnchor.ReplaceAllStringFunc(html, func(tag string) string {
			if !strings.Contains(tag, `href="http`) || strings.Contains(tag, "target=") {
				return tag
			}
			return strings.TrimSuffix(tag, ">") + ` target="_blank" rel="noopener noreferrer">`
		})

		_, err := io.WriteString(w, html)
		return err
	})
}
