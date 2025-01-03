package meta

import (
	"strings"
	"testing"
)

func TestExtractMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		content  string
		path     string
		expected Meta
		found    bool
	}{
		{
			name:    "valid templ-static directive with path",
			content: "//templ-static path:/static/assets\n\npackage main",
			path:    "/some/path/file.go",
			expected: Meta{
				FilePath: "/some/path/file.go",
				Path:     "/static/assets",
			},
			found: true,
		},
		{
			name:    "valid templ-static directive without path",
			content: "//templ-static\n\npackage main",
			path:    "/some/path/file.go",
			expected: Meta{
				FilePath: "/some/path/file.go",
				Path:     "",
			},
			found: true,
		},
		{
			name:     "no templ-static directive",
			content:  "package main",
			path:     "/some/path/file.go",
			expected: Meta{},
			found:    false,
		},
		{
			name:     "templ-static directive after package",
			content:  "package main\n//templ-static path:/static/assets",
			path:     "/some/path/file.go",
			expected: Meta{},
			found:    false,
		},
		{
			name:     "misformatted: templ-static with extra spaces",
			content:  "//    templ-static    path:  /static/assets\npackage main",
			path:     "/some/path/file.go",
			expected: Meta{},
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader(tt.content)
			meta, found := Extract(reader, tt.path)

			if found != tt.found {
				t.Errorf("expected found = %v, got %v", tt.found, found)
			}

			if meta != tt.expected {
				t.Errorf("expected meta = %+v, got %+v", tt.expected, meta)
			}
		})
	}
}
