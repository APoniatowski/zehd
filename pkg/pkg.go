package pkg

import (
	"fmt"
	"net/http"
	"zehd/pkg/env"
)

// Common variables
var (
	RefreshCache    string
	CollectionError error
	TemplatesDir    string
	TemplateType    string
	Javascript      string
	CSS             string
	Images          string
	Downloads       string
	GitLink         string
	GitUsername     string
	GitToken        string
)

const Disable = "unused-path" // TODO: might change this later

func init() {
	fmt.Print("Loading Configuration via Environment Variables...")
	config := []struct {
		value    *string
		getValue func() (string, int)
	}{
		{&TemplatesDir, env.EnvTemplateDir},
		{&TemplateType, env.EnvTemplateType},
		{&RefreshCache, env.EnvCacheRefresh},
		{&Javascript, func() (string, int) { return env.EnvPathDefiner("js") }},
		{&CSS, func() (string, int) { return env.EnvPathDefiner("css") }},
		{&Images, func() (string, int) { return env.EnvPathDefiner("images") }},
		{&Downloads, func() (string, int) { return env.EnvPathDefiner("downloads") }},
		{&GitLink, env.EnvGitLink},
		{&GitUsername, env.EnvGitUsername},
		{&GitToken, env.EnvGitToken},
	}
	configScoreTotal := 0
	score := 0
	for _, c := range config {
		*c.value, score = c.getValue()
		configScoreTotal += score
	}
	switch {
	case configScoreTotal == 20:
		fmt.Println(" Fully configured.")
	case configScoreTotal > 12:
		fmt.Println("Mostly configured, defaults set on some settings.")
	case configScoreTotal > 6:
		fmt.Println(" Partially configured, defaults set on some settings.")
	default:
		fmt.Println(" Using defaults.")
	}
}

// Common functions
// faviconHandler the gremlin at the top of your browser tab.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, TemplatesDir+"favicon/favicon.png")
}
