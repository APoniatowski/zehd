package caching

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zehd/pkg"

	"github.com/APoniatowski/boillog"
)

// CachePages Method that walks the specified or default directories and caches the templates
func (pages *Pages) CachePages() error {
	defer boillog.TrackTime("cacher", time.Now())

	if len(pkg.GitLink) != 0 {
		err := Git("refresh")
		if err != nil {
			boillog.LogIt("CachePages", "ERROR", err.Error())
			return err
		}

		return nil
	}

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

		case ".org":
			filetype = ".org"

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
