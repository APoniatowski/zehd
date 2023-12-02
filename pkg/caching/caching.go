package caching

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"zehd-frontend/pkg"
	"strings"
	"time"

	"github.com/APoniatowski/boillog"
)

// CachePages Method that walks the specified or default directories and caches the templates
func (pages *Pages) CachePages() error {
	defer boillog.TrackTime("cacher", time.Now())
	errchdir := os.Chdir(pkg.TemplatesDir + pkg.TemplateType)
	if errchdir != nil {
		boillog.LogIt("cachepages", "error", "chdir returned an error: "+fmt.Sprintln(errchdir))
	}
	err := filepath.WalkDir(pkg.TemplatesDir+pkg.TemplateType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		pathtoremove := pkg.TemplatesDir + pkg.TemplateType
		indextoremove := strings.Index(path, pathtoremove)
		if indextoremove == -1 {
			boillog.LogIt("cachepages", "error", "directory not found")
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
		tmpl, err := templateBuilder(croppedtemplatepath, filetype)
		if err != nil {
			return fmt.Errorf("failed to build template for file %q: %v", path, err)
		}
		pages.RouteMap[name] = tmpl
		return nil
	})
	if err != nil {
		boillog.LogIt("cachepages", "error", "walkdir returned an error: "+fmt.Sprintln(err))
		return err
	}
	return nil
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
