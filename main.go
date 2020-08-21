package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/shinshin86/metatag-gen/templates"
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

	case "pug":
		tmplType = "pug.tmpl"

	case "haml":
		tmplType = "haml.tmpl"

	case "slim":
		tmplType = "slim.tmpl"

	case "jsx":
		tmplType = "jsx.tmpl"

	default:
		log.Fatal("Invalid template type")
	}

	f, err := templates.Root.Open(path.Join("/", tmplType))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New(tmplType).Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}

	if err = t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
