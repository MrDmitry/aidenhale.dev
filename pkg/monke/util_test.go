package monke

import (
	"testing"
)

func TestSanitizeUrl(t *testing.T) {
	testData := map[string]string{
		"":                             "",                             // empty string stays empty
		"a":                            "/a",                           // path becomes absolute
		"a/":                           "/a/",                          // path becomes absolute
		"/":                            "/",                            // root stays root
		"/root":                        "/root",                        // path stays path
		"/root/":                       "/root/",                       // trailing slash is preserved
		"/root/../other":               "/other",                       // relative paths are resolved
		"/root/../../other":            "/other",                       // relative paths are resolved
		"/🤡":                           "/%F0%9F%A4%A1",                // funny symbols are encoded
		"/%F0%9F%A4%A1":                "/%F0%9F%A4%A1",                // encoded paths are preserved
		"/path?query":                  "/path?query",                  // query is preserved
		"/path/?query":                 "/path/?query",                 // query is preserved
		"/path/?param=/path/":          "/path/?param=/path/",          // query is preserved
		"/path/?param=/path/../../../": "/path/?param=/path/../../../", // query is preserved
		"https://bad.url//":            "https://bad.url/",             // works with URLs
		"https://bad.url/root/.":       "https://bad.url/root",         // works with URLs but resolves . funny
		"https://bad.url/root/./":      "https://bad.url/root/",        // recognizes a trailing slash
		"https://bad.url/root/../../":  "https://bad.url/",             // respects path boundary
	}
	for s, want := range testData {
		got := SanitizeUrl(s)
		if got != want {
			t.Fatalf("Failed to sanitize %s: expected %s, got %s", s, want, got)
		}
	}
}
