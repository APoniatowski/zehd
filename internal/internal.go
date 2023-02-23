package internal

import (
	"net/http"
	"poniatowski-dev/internal/env"
)

// Common variables
var CollectionError error
var TemplatesDir = env.EnvTemplateDir()
var TemplateType = env.EnvTemplateType()

// Common functions
// faviconHandler the gremlin at the top of your browser tab.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, TemplatesDir+"favicon/favicon.png")
}
