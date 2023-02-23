package caching

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"poniatowski-dev/internal"
	"poniatowski-dev/internal/logging"
)

type Pages struct {
	RouteMap map[string]*template.Template
}

func (pages *Pages) CachePages() error {
	errchdir := os.Chdir(internal.TemplatesDir + internal.TemplateType)
	if errchdir != nil {
		logging.LogIt("cachepages", "error", "chdir returned an error: "+fmt.Sprintln(errchdir))
	}
	err := filepath.WalkDir(internal.TemplatesDir+internal.TemplateType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		pathtoremove := internal.TemplatesDir + internal.TemplateType
		indextoremove := strings.Index(path, pathtoremove)
		if indextoremove == -1 {
			logging.LogIt("cachepages", "error", "directory not found")
		}
		croppedtemplatepath := strings.TrimPrefix(path[indextoremove+len(pathtoremove):], "/")
		var filetype string
		switch filepath.Ext(path) {
		case ".gohtml":
			filetype = ".gohtml"
		case ".html":
			filetype = ".html"
		case ".md":
			filetype = ".md"
		default:
			filetype = "invalid"
		}
		name := strings.TrimSuffix(croppedtemplatepath, filepath.Ext(path))
		tmpl, err := templatebuilder(croppedtemplatepath, filetype)
		if err != nil {
			return fmt.Errorf("failed to build template for file %q: %v", path, err)
		}
		pages.RouteMap[name] = tmpl
		return nil
	})
	if err != nil {
		logging.LogIt("cachepages", "error", "walkdir returned an error: "+fmt.Sprintln(err))
		return err
	}
	return nil
}

func templatebuilder(page, filetype string) (*template.Template, error) {
	if filetype == "invalid" {
		return nil, errors.New("invalid filetype: " + page)
	}
	layoutpage := internal.TemplatesDir + "layout." + internal.TemplateType
	templatepage := internal.TemplatesDir + internal.TemplateType + "/" + page
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
