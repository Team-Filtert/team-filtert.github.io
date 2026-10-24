package templates

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/LukasDerBaum42/astwerk/ssg"
	templpkg "github.com/a-h/templ"
)

// MarkdownBody renders a markdown body with the page's Ctx, so
// content.Page.Render can resolve root-relative URLs and run components.
// External links are handled by the loader's render options.
func MarkdownBody(c ssg.Ctx, body func(ssg.Ctx) templpkg.Component) templpkg.Component {
	return templpkg.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if body == nil {
			return nil
		}
		component := body(c)
		if component == nil {
			return nil
		}
		return component.Render(ctx, w)
	})
}

// MemberBio renders a member's bio, with a visible fallback when it is empty.
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

		if strings.TrimSpace(buf.String()) == "" {
			_, err := io.WriteString(w, `<p class="no-bio">No bio available.</p>`)
			return err
		}

		_, err := w.Write(buf.Bytes())
		return err
	})
}
