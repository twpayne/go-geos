package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"regexp"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

var (
	templateDataFilename = flag.String("data", "", "data filename")
	outputFilename       = flag.String("output", "", "output filename")
)

func run() error {
	flag.Parse()

	var templateData any
	if *templateDataFilename != "" {
		dataBytes, err := os.ReadFile(*templateDataFilename)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(dataBytes, &templateData); err != nil {
			return err
		}
	}

	if flag.NArg() == 0 {
		return fmt.Errorf("no arguments")
	}

	templateName := path.Base(flag.Arg(0))
	buffer := &bytes.Buffer{}
	funcMap := template.FuncMap{
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
		//nolint:forbidigo
		fmt.Println(err)
		os.Exit(1)
	}
}
