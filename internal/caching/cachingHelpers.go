package caching

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/russross/blackfriday/v2"

	"zehd-frontend/internal/logging"
)

func pageBuilder(templatePath, layoutPath string) (*template.Template, error) {
	defer logging.TrackTime("page-builder", time.Now())
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
	defer logging.TrackTime("org-converter", time.Now())
	templates := template.New("")

	_, notFoundErr := os.Stat(layoutPath)
	if notFoundErr != nil {
		orgBytes, err := os.ReadFile(orgPath)
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
		orgBytes, err := os.ReadFile(orgPath)
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
	defer logging.TrackTime("md-converter", time.Now())
	templates := template.New("")
	_, notFoundErr := os.Stat(layoutPath)
	if notFoundErr != nil {
		markdownBytes, err := os.ReadFile(layoutPath)
		if err != nil {
			return nil, err
		}
		html := string(markdown.ToHTML(markdownBytes, nil, nil))
		html = fmt.Sprintf(`{{define "markdown"}}%s{{end}}`, html)
		templates, err = templates.Parse(html)
		if err != nil {
			return nil, err
		}
	} else {
		markdownBytes, err := os.ReadFile(markdownPath)
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
