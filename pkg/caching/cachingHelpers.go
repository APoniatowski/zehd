package caching

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"
	"zehd/pkg"

	"github.com/APoniatowski/boillog"
	"github.com/APoniatowski/funcmytemplate"

	"github.com/gomarkdown/markdown"
	"github.com/russross/blackfriday/v2"
)

// pageBuilder Private helper function that builds HTML/goHTML pages and returns the templates
func pageBuilder(templatePath, layoutPath string) (*template.Template, error) {
	defer boillog.TrackTime("page-builder", time.Now())
	funcmytemplates := funcmytemplate.Add()
	templates := template.New("layout")
	_, err := os.Stat(layoutPath)
	if err != nil {
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
	var errReturn error
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
		if err != nil {
			errReturn = err
		}
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
		if err != nil {
			errReturn = err
		}
	}
	if errReturn != nil {
		return nil, errReturn
	}
	return templates, nil
}

// convertMarkdownToTemplate Private helper function that builds markdown pages and returns the templates
func convertMarkdownToTemplate(markdownPath, layoutPath string) (*template.Template, error) {
	defer boillog.TrackTime("md-converter", time.Now())
	var errReturn error
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
		if err != nil {
			errReturn = err
		}
		templates, err = templates.Funcs(funcmytemplates).Parse(html)
		if err != nil {
			errReturn = err
		}
	}
	if errReturn != nil {
		return nil, errReturn
	}
	return templates, nil
}

// templateBuilder Private function for building templates, which is called by CachePages
func templateBuilder(page, filetype string) (*template.Template, error) {
	defer boillog.TrackTime("template-builder", time.Now())
	if filetype == "invalid" {
		return nil, errors.New("invalid filetype: " + page)
	}
	layoutpage := pkg.TemplatesDir + "layout." + pkg.TemplateType
	templatepage := pkg.TemplatesDir + pkg.TemplateType + "/" + page
	_, notfounderr := os.Stat(templatepage)
	if notfounderr != nil {
		if os.IsNotExist(notfounderr) {
			return nil, errors.New("template does not exist: " + page)
		}
	}
	var parseerr error
	var templates *template.Template
	switch filetype {
	case ".org":
		templates, parseerr = convertOrgToTemplate(templatepage, layoutpage)
	case ".md":
		templates, parseerr = convertMarkdownToTemplate(templatepage, layoutpage)
	default:
		templates, parseerr = pageBuilder(templatepage, layoutpage)
	}
	if parseerr != nil {
		return nil, errors.New("error parsing templates: " + fmt.Sprintln(parseerr))
	}
	return templates, nil
}
