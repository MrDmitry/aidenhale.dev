package monke

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/gommon/log"
)

type renderHookData struct {
	prefix string
}

// `hx-boost` updates the history only after a request to an <img src="./..."/> fails with 404 (requesting an asset
// from the old URL path) but I want to keep using `hx-boost` so I gotta hack it
func makeRenderHook(data *renderHookData) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if image, ok := node.(*ast.Image); ok && image.Destination[0] == '.' {
			relPath := path.Join(data.prefix, string(image.Destination))
			result, err := filepath.Abs(relPath)

			if err == nil {
				image.Destination = []byte(result)
			} else {
				log.Warnf("failed to resolve absolute path for %s: %+v", relPath, err)
			}
		}
		return ast.GoToNext, false
	}
}

func mdToHTML(md []byte, absPrefix string) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: makeRenderHook(&renderHookData{prefix: absPrefix})}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func RenderMarkdownAbs(f string, absPrefix string) ([]byte, error) {
	file, err := os.ReadFile(f)

	if err != nil {
		return nil, err
	}

	return mdToHTML(file, absPrefix), nil
}

func RenderMarkdown(f string) ([]byte, error) {
	return RenderMarkdownAbs(f, "")
}
