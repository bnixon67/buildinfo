// Package buildinfo provides functions to format and retrieve Go build
// information in both text and JSON formats.
package buildinfo

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

// QuoteIf returns the given string s wrapped in quotes if quote is true.
// If quote is false, it returns s unchanged.
func quoteIf(s string, quote bool) string {
	if quote {
		return strconv.Quote(s)
	}
	return s
}

// formatModule returns Module as formated text.
func formatModule(sb *strings.Builder, module debug.Module, quote bool, prefix, indent string) {
	fmt.Fprintf(sb, "%s%sPath: %s\n", prefix, indent, quoteIf(module.Path, quote))
	fmt.Fprintf(sb, "%s%sVersion: %s\n", prefix, indent, quoteIf(module.Version, quote))
	fmt.Fprintf(sb, "%s%sSum: %s\n", prefix, indent, quoteIf(module.Sum, quote))
	if module.Replace != nil {
		fmt.Fprintln(sb, prefix+indent+"Replace:")
		formatModule(sb, *module.Replace, quote, prefix+indent, indent)
	}
}

// FormatText returns BuildInfo as formatted text.
func FormatText(bi debug.BuildInfo, quote bool, prefix, indent string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%sGo Version: %s\n", prefix, quoteIf(bi.GoVersion, quote))
	fmt.Fprintf(&sb, "%sPath: %s\n", prefix, quoteIf(bi.Path, quote))

	fmt.Fprintf(&sb, "%sMain Module:\n", prefix)
	formatModule(&sb, bi.Main, quote, prefix, indent)

	if len(bi.Deps) > 0 {
		fmt.Fprintf(&sb, "%sDependencies:\n", prefix)
		for _, dep := range bi.Deps {
			formatModule(&sb, *dep, quote, prefix, indent)
		}
	}

	if len(bi.Settings) > 0 {
		fmt.Fprintf(&sb, "%sSettings:\n", prefix)
		for _, setting := range bi.Settings {
			fmt.Fprintf(&sb, "%s%s%s: %s\n", prefix, indent, setting.Key, quoteIf(setting.Value, quote))
		}
	}

	return sb.String()
}

// FormatJSON returns the BuildInfo as a JSON string.
func FormatJSON(info debug.BuildInfo, prefix, indent string) (string, error) {
	jsonData, err := json.MarshalIndent(info, prefix, indent)
	if err != nil {
		return "", fmt.Errorf("cannot convert build info to JSON: %w", err)
	}
	return prefix + string(jsonData), nil
}
