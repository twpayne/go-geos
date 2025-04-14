// execute-template executes a Go template with data.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"regexp"
	"slices"
	"text/template"
	"unicode"

	"github.com/goccy/go-yaml"
)

var (
	templateDataFilename = flag.String("data", "", "data filename")
	outputFilename       = flag.String("output", "", "output filename")

	cTypes = map[string]string{
		"BufCapStyle":            "C.int",
		"BufJoinStyle":           "C.int",
		"PrecisionRule":          "C.int",
		"RelateBoundaryNodeRule": "C.int",
		"float64":                "C.double",
		"int":                    "C.int",
		"uint":                   "C.unsigned",
	}
)

func run() error {
	flag.Parse()

	var templateData []map[string]any
	if *templateDataFilename != "" {
		dataBytes, err := os.ReadFile(*templateDataFilename)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(dataBytes, &templateData); err != nil {
			return err
		}
	}

	if !slices.IsSortedFunc(templateData, func(a, b map[string]any) int {
		switch aName, bName := a["name"].(string), b["name"].(string); { //nolint:forcetypeassert
		case aName < bName:
			return -1
		case aName == bName:
			return 0
		default:
			return 1
		}
	}) {
		return errors.New("template data not sorted by name")
	}

	if flag.NArg() == 0 {
		return errors.New("no arguments")
	}

	templateName := path.Base(flag.Arg(0))
	buffer := &bytes.Buffer{}
	funcMap := template.FuncMap{
		"cType": func(goType string) string {
			cType, ok := cTypes[goType]
			if !ok {
				panic(errors.New(goType + ": unknown C type for Go type"))
			}
			return cType
		},
		"fatal": func(s string) string {
			panic(s)
		},
		"firstRuneToLower": func(s string) string {
			runes := []rune(s)
			runes[0] = unicode.ToLower(runes[0])
			return string(runes)
		},
		"replaceAllRegexp": func(expr, repl, s string) string {
			return regexp.MustCompile(expr).ReplaceAllString(s, repl)
		},
	}
	tmpl, err := template.New(templateName).Funcs(funcMap).ParseFiles(flag.Args()...)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(buffer, templateData); err != nil {
		return err
	}

	output, err := format.Source(buffer.Bytes())
	if err != nil {
		output = buffer.Bytes()
	}

	if *outputFilename == "" {
		if _, err := os.Stdout.Write(output); err != nil {
			return err
		}
	} else if data, err := os.ReadFile(*outputFilename); err != nil || !bytes.Equal(data, output) {
		//nolint:gosec
		if err := os.WriteFile(*outputFilename, output, 0o666); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
