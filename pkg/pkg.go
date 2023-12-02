package pkg

import (
	"net/http"
	"zehd-frontend/pkg/env"
)

// Common variables
var (
	CollectionError error
	TemplatesDir    = env.EnvTemplateDir()
	TemplateType    = env.EnvTemplateType()
)

// Common functions
// faviconHandler the gremlin at the top of your browser tab.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, TemplatesDir+"favicon/favicon.png")
}
