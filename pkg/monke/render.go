package monke

import (
	"io"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type renderHookData struct {
	prefix []byte
}

// html.AbsolutePrefix doesn't work for links that start with ./ while `hx-boost` updates the history only after a
// request to an <img src="./..."/> fails with 404 (requesting an asset from the old URL path)
func makeRenderHook(data *renderHookData) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if image, ok := node.(*ast.Image); ok && image.Destination[0] == '.' {
			image.Destination = append(data.prefix, image.Destination[2:]...)
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
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: makeRenderHook(&renderHookData{prefix: []byte(absPrefix)})}
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
