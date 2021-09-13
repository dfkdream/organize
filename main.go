package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type file struct {
	From string
	To   string
	Info os.FileInfo
}

func buildFileList(root, pattern string, recursive bool, outputDir string, outputTemplate *template.Template) ([]file, error) {
	fileList := make([]file, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}

		if info.IsDir() {
			if recursive || path == root {
				return nil
			}

			return filepath.SkipDir
		}

		match, mErr := filepath.Match(pattern, info.Name())
		if mErr != nil {
			return mErr
		}

		if match {
			f := file{
				From: path,
				Info: info,
			}

			buff := new(strings.Builder)
			tErr := outputTemplate.Execute(buff, f)
			if tErr != nil {
				return tErr
			}

			f.To = filepath.Join(outputDir, buff.String())

			fileList = append(fileList, f)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func main() {
	inputDir := flag.String("i", ".", "Input Directory")
	outputDir := flag.String("o", "./output", "Output Directory")
	pattern := flag.String("p", "*", "File name pattern (Default: *)")
	recursive := flag.Bool("r", false, "Enable recursive walk")
	dryRun := flag.Bool("dry-run", false, "Enable dry run")
	outputTemplateString := flag.String("t", "{{ .Info.Name }}", "Output file template")
	flag.Parse()

	fmt.Println("Input Directory: ", *inputDir)
	fmt.Println("Input Pattern:   ", *pattern)
	fmt.Println("Output Directory:", *outputDir)
	fmt.Println("Output Template: ", *outputTemplateString)
	fmt.Println()

	if *dryRun {
		fmt.Println("*** DRY RUN ***")
		fmt.Println()
	}

	outputTemplate, err := template.New("output").
		Funcs(template.FuncMap{
			"count": func() func(string) string {
				i := 0
				return func(format string) string {
					i = i + 1
					return fmt.Sprintf(format, i)
				}
			}(),
			"ext": filepath.Ext,
			"chTimes": func(filename string, t time.Time) (string, error) {
				fmt.Println("Change file", filename, "timestamp to", t)
				if *dryRun {
					return "", nil
				}
				return "", os.Chtimes(filename, t, t)
			},
			"parseTime": time.Parse,
		}).
		Parse(*outputTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	fl, err := buildFileList(*inputDir, *pattern, *recursive, *outputDir, outputTemplate)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fl {
		dir := filepath.Dir(f.To)
		if dir != *outputDir {
			fmt.Printf("mkdir %s (0777)\n", filepath.Dir(f.To))
			if !*dryRun {
				err = os.MkdirAll(filepath.Dir(f.To), 0777)
				if err != nil {
					log.Println(err)
				}
			}
		}

		fmt.Printf("mv %s -> %s\n", f.From, f.To)

		if !*dryRun {

			err = os.Rename(f.From, f.To)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
