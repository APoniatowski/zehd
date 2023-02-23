package caching

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/russross/blackfriday/v2"
)

func pageBuilder(templatePath, layoutPath string) (*template.Template, error) {
	templates := template.New("")
	_, err := os.Stat(layoutPath)
	if err != nil {
		templates, err = template.ParseFiles(templatePath)
		if err != nil {
			return nil, err
		}
	} else {
		templates, err = template.ParseFiles(layoutPath, templatePath)
		if err != nil {
			return nil, err
		}
	}
	return templates, nil
}

func convertOrgToTemplate(orgPath, layoutPath string) (*template.Template, error) {
	templates := template.New("")

	_, notFoundErr := os.Stat(layoutPath)
	if notFoundErr != nil {
		orgBytes, err := ioutil.ReadFile(orgPath)
		if err != nil {
			return nil, err
		}

		htmlBytes := blackfriday.Run(orgBytes)
		html := string(htmlBytes)

		html = fmt.Sprintf(`{{define "org"}}%s{{end}}`, html)
		templates, err = templates.Parse(html)
		if err != nil {
			return nil, err
		}
	} else {
		orgBytes, err := ioutil.ReadFile(orgPath)
		if err != nil {
			return nil, err
		}

		htmlBytes := blackfriday.Run(orgBytes)
		html := string(htmlBytes)

		html = fmt.Sprintf(`{{define "org"}}%s{{end}}`, html)

		templates, err = template.ParseFiles(layoutPath)
		if err != nil {
			return nil, err
		}

		templates, err = templates.Parse(html)
		if err != nil {
			return nil, err
		}
	}
	return templates, nil
}

func convertMarkdownToTemplate(markdownPath, layoutPath string) (*template.Template, error) {
	templates := template.New("")
	_, notFoundErr := os.Stat(layoutPath)
	if notFoundErr != nil {
		markdownBytes, err := ioutil.ReadFile(layoutPath)
		if err != nil {
			return nil, err
		}
		html := string(markdown.ToHTML(markdownBytes, nil, nil))
		html = fmt.Sprintf(`{{define "markdown"}}%s{{end}}`, html)
		templates, err = templates.Parse(html)
	} else {
		markdownBytes, err := ioutil.ReadFile(markdownPath)
		if err != nil {
			return nil, err
		}

		html := string(markdown.ToHTML(markdownBytes, nil, nil))
		html = fmt.Sprintf(`{{define "markdown"}}%s{{end}}`, html)

		templates, err = template.ParseFiles(layoutPath)
		if err != nil {
			return nil, err
		}
		templates, err = templates.Parse(html)
		if err != nil {
			return nil, err
		}
	}
	return templates, nil
}
