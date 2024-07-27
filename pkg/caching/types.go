package caching

import (
	"html/template"
)

// Pages Struct for caching templates and routes
type Pages struct {
	RouteMap map[string]*template.Template
}
