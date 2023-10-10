package monke

import (
	"errors"
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

type customHTMLRenderHookData struct {
	prefix string
}

// `hx-boost` updates the history only after a request to an <img src="./..."/> fails with 404 (requesting an asset
// from the old URL path) but I want to keep using `hx-boost` so I gotta hack it
func customHTMLRenderHook(data *customHTMLRenderHookData) html.RenderNodeFunc {
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
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: customHTMLRenderHook(&customHTMLRenderHookData{prefix: absPrefix}),
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func RenderMarkdownToHTMLAbs(f string, absPrefix string) ([]byte, error) {
	if f == "" {
		return nil, errors.New("empty filepath given")
	}

	file, err := os.ReadFile(f)

	if err != nil {
		return nil, err
	}

	return mdToHTML(file, absPrefix), nil
}

func RenderMarkdownToHTML(f string) ([]byte, error) {
	return RenderMarkdownToHTMLAbs(f, "")
}

type previewTracker struct {
	written     int
	threshold   int
	wantSpace   bool
	terminating bool
}

func previewRenderer(data *previewTracker) html.RenderNodeFunc {
	data.written = 0
	data.wantSpace = false
	data.terminating = false

	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if data.written >= data.threshold {
			if !data.terminating {
				w.Write([]byte("..."))
			}
			data.terminating = true
			return ast.Terminate, true
		}
		switch node := node.(type) {
		case *ast.Text:
			size := min(len(node.Literal), data.threshold-data.written)
			w.Write(node.Literal[:size])
			data.wantSpace = (node.Literal[size-1] != ' ')
			data.written += size
		case *ast.Paragraph:
			// Separate paragraphs with spaces
			if entering && data.wantSpace {
				w.Write([]byte(" "))
				data.written += 1
			}
			break
		case *ast.Document:
			break
		case *ast.Emph:
			break
		case *ast.Strong:
			break
		case *ast.BlockQuote:
			break
		case *ast.Link:
			break
		case *ast.Citation:
			break
		case *ast.Subscript:
			break
		case *ast.Superscript:
			break
		default:
			return ast.SkipChildren, true
		}
		return ast.GoToNext, true
	}
}

func mdToTextPreview(md []byte, t int) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: previewRenderer(&previewTracker{threshold: t}),
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func RenderMarkdownToTextPreview(f string, t int) ([]byte, error) {
	file, err := os.ReadFile(f)

	if err != nil {
		return nil, err
	}

	return mdToTextPreview(file, t), nil
}
