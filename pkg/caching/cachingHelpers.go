package caching

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/APoniatowski/boillog"
	"github.com/APoniatowski/funcmytemplate"

	"github.com/gomarkdown/markdown"
	"github.com/russross/blackfriday/v2"
)

// pageBuilder Private helper function that builds HTML/goHTML pages and returns the templates
func pageBuilder(templatePath, layoutPath string) (*template.Template, error) {
	defer boillog.TrackTime("page-builder", time.Now())
	funcmytemplates := funcmytemplate.Add()
	templates := template.New("")
	_, err := os.Stat(layoutPath)
	if err != nil || !os.IsNotExist(err) {
		templates, err = templates.Funcs(funcmytemplates).ParseFiles(templatePath)
	} else {
		templates, err = templates.Funcs(funcmytemplates).ParseFiles(layoutPath, templatePath)
	}
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// convertOrgToTemplate Private helper function that builds org-mode pages and returns the templates
func convertOrgToTemplate(orgPath, layoutPath string) (*template.Template, error) {
	defer boillog.TrackTime("org-converter", time.Now())
	var err error
	funcmytemplates := funcmytemplate.Add()
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
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
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
		templates, err = templates.Funcs(funcmytemplates).ParseFiles(layoutPath)
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
	}
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// convertMarkdownToTemplate Private helper function that builds markdown pages and returns the templates
func convertMarkdownToTemplate(markdownPath, layoutPath string) (*template.Template, error) {
	defer boillog.TrackTime("md-converter", time.Now())
	var err error
	funcmytemplates := funcmytemplate.Add()
	templates := template.New("")
	_, notFoundErr := os.Stat(layoutPath)
	if notFoundErr != nil {
		markdownBytes, err := os.ReadFile(layoutPath)
		if err != nil {
			return nil, err
		}
		html := string(markdown.ToHTML(markdownBytes, nil, nil))
		html = fmt.Sprintf(`{{define "markdown"}}%s{{end}}`, html)
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
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
		templates, err = templates.Funcs(funcmytemplates).ParseFiles(layoutPath)
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
	}
	if err != nil {
		return nil, err
	}
	return templates, nil
}
