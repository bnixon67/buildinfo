package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/bnixon67/buildinfo"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	prefix := flag.String("prefix", "", "Prefix for each line in text format")
	indent := flag.String("indent", "  ", "Indentation string for nested structures in text format")
	quote := flag.Bool("quote", true, "Quote values")
	flag.Parse()

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("No build info available.")
		return
	}

	switch *format {
	case "text":
		fmt.Println(buildinfo.FormatText(*info, *quote, *prefix, *indent))
	case "json":
		jsonData, err := buildinfo.FormatJSON(*info, *prefix, *indent)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(jsonData)
	default:
		fmt.Printf("Invalid format: %q. Use 'text' or 'json'.\n", *format)
		os.Exit(1)
	}
}
