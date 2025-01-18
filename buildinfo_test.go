package buildinfo

import (
	"encoding/json"
	"runtime/debug"
	"strings"
	"testing"
)

// TestQuoteIf tests the quoteIf function using a table-driven approach.
func TestQuoteIf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		quote bool
		want  string
	}{
		{"No Quoting", "hello", false, "hello"},
		{"With Quoting", "hello", true, `"hello"`},
		{"Empty String Without Quoting", "", false, ""},
		{"Empty String With Quoting", "", true, "\"\""},
		{"String With Newline", "hello\nworld", true, `"hello\nworld"`},
		{"String With Quotes", `"test"`, true, `"\"test\""`},
		{"Numerical String Without Quoting", "12345", false, "12345"},
		{"Numerical String With Quoting", "12345", true, `"12345"`},
		{"Special Characters", "Go: \tGopher", true, `"Go: \tGopher"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := quoteIf(tt.input, tt.quote)
			if result != tt.want {
				t.Errorf("quoteIf(%q, %v)\ngot  %q\nwant %q",
					tt.input, tt.quote, result, tt.want)
			}
		})
	}
}

// TestFormatText tests the FormatText function using a table-driven approach.
func TestFormatText(t *testing.T) {
	tests := []struct {
		name   string
		bi     debug.BuildInfo
		quote  bool
		prefix string
		indent string
		want   string
	}{
		{
			name: "Minimal BuildInfo",
			bi: debug.BuildInfo{
				GoVersion: "go1.20.3",
				Path:      "example.com/project",
				Main: debug.Module{
					Path:    "example.com/project",
					Version: "v1.2.3",
					Sum:     "h1:abcdef1234567890",
				},
			},
			quote:  false,
			prefix: "",
			indent: "  ",
			want: `Go Version: go1.20.3
Path: example.com/project
Main Module:
  Path: example.com/project
  Version: v1.2.3
  Sum: h1:abcdef1234567890
`,
		},
		{
			name: "With Replace",
			bi: debug.BuildInfo{
				GoVersion: "go1.20.3",
				Path:      "example.com/project",
				Main: debug.Module{
					Path:    "example.com/project",
					Version: "v1.2.3",
					Sum:     "h1:abcdef1234567890",
					Replace: &debug.Module{
						Path:    "example.com/replace",
						Version: "v4.5.6",
					},
				},
			},
			quote:  false,
			prefix: "",
			indent: "  ",
			want: `Go Version: go1.20.3
Path: example.com/project
Main Module:
  Path: example.com/project
  Version: v1.2.3
  Sum: h1:abcdef1234567890
  Replace:
    Path: example.com/replace
    Version: v4.5.6
    Sum: 
`,
		},
		{
			name: "With Dependencies",
			bi: debug.BuildInfo{
				GoVersion: "go1.19",
				Path:      "github.com/test/project",
				Main: debug.Module{
					Path:    "github.com/test/project",
					Version: "v2.0.0",
					Sum:     "h1:def4567890abcdef",
				},
				Deps: []*debug.Module{
					{
						Path:    "github.com/dependency/module",
						Version: "v1.5.0",
						Sum:     "h1:123abc456def7890",
					},
				},
			},
			quote:  false,
			prefix: "",
			indent: "  ",
			want: `Go Version: go1.19
Path: github.com/test/project
Main Module:
  Path: github.com/test/project
  Version: v2.0.0
  Sum: h1:def4567890abcdef
Dependencies:
  Path: github.com/dependency/module
  Version: v1.5.0
  Sum: h1:123abc456def7890
`,
		},
		{
			name: "With Settings and Quoting Enabled",
			bi: debug.BuildInfo{
				GoVersion: "go1.18",
				Path:      "project.com/foo",
				Main: debug.Module{
					Path:    "project.com/foo",
					Version: "v1.0.1",
					Sum:     "h1:xyz9876543210",
				},
				Settings: []debug.BuildSetting{
					{Key: "CGO_ENABLED", Value: "1"},
					{Key: "GOOS", Value: "linux"},
				},
			},
			quote:  true,
			prefix: "--> ",
			indent: "    ",
			want: `--> Go Version: "go1.18"
--> Path: "project.com/foo"
--> Main Module:
-->     Path: "project.com/foo"
-->     Version: "v1.0.1"
-->     Sum: "h1:xyz9876543210"
--> Settings:
-->     CGO_ENABLED: "1"
-->     GOOS: "linux"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatText(tt.bi, tt.quote, tt.prefix, tt.indent)
			if result != tt.want {
				t.Errorf("\n--- Got:\n%s\n--- Want:\n%s", result, tt.want)
			}
		})
	}
}

// TestFormatJSON tests the FormatJSON function using a table-driven approach.
func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		bi      debug.BuildInfo
		prefix  string
		indent  string
		wantErr bool
	}{
		{
			name: "Basic JSON output",
			bi: debug.BuildInfo{
				GoVersion: "go1.21",
				Path:      "test.com/pkg",
				Main: debug.Module{
					Path:    "test.com/pkg",
					Version: "v0.9.1",
					Sum:     "h1:abc987654321",
				},
			},
			prefix:  "",
			indent:  "  ",
			wantErr: false,
		},
		{
			name:    "Empty BuildInfo",
			bi:      debug.BuildInfo{},
			prefix:  "",
			indent:  "  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatJSON(tt.bi, tt.prefix, tt.indent)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify that result is valid JSON
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(strings.TrimPrefix(result, tt.prefix)), &parsed); err != nil {
				t.Errorf("Invalid JSON output: %v", err)
			}
		})
	}
}
