package main

import (
	"flag"
	"log"
	"os"
	"path"
	"text/template"
)

func main() {
	var (
		url         = flag.String("u", "", "URL")
		title       = flag.String("t", "", "Title")
		description = flag.String("d", "", "Description")
		keywords    = flag.String("k", "", "Keywords (If you want to specify more than one, please separate them with a comma.)")
		imgPath     = flag.String("i", "", "OGP image path")
		tmpl        = flag.String("tmpl", "html", "Use template")
	)
	flag.Parse()

	data := map[string]string{
		"Url":         *url,
		"Title":       *title,
		"Description": *description,
		"Keywords":    *keywords,
		"ImgPath":     *imgPath,
	}

	var tmplType string

	switch *tmpl {
	case "html":
		tmplType = "html.tmpl"

	default:
		log.Fatal("Invalid template type")
	}

	t, err := template.New(tmplType).ParseFiles(path.Join("template", tmplType))
	if err != nil {
		log.Fatal(err)
	}

	if err = t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
