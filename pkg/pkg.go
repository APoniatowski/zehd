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
	fmt.Println("Loading Configuration via Environment Variables.")

	orange := "\033[48;5;208m" // Orange background
	green := "\033[42m"        // Green background
	reset := "\033[0m"         // Reset to default

	orangeBlock := orange + "  " + reset
	greenBlock := green + "  " + reset

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
		switch score {
		case 1:
			fmt.Print(orangeBlock)
		case 2:
			fmt.Print(greenBlock)
		}
		configScoreTotal += score
	}

	fmt.Println("")

	switch {
	case configScoreTotal == 20:
		fmt.Println("      [--- Fully configured. ---]")

	case configScoreTotal > 15:
		fmt.Println("      [--- Mostly configured, defaults set on some settings. ---]")

	case configScoreTotal > 10:
		fmt.Println("      [--- Partially configured, defaults set on most settings. ---]")

	case configScoreTotal == len(config):
		fmt.Println("      [--- Using defaults. ---]")
	}

	fmt.Println("<======================================================================================>")
}

// Common functions
// faviconHandler the gremlin at the top of your browser tab.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, TemplatesDir+"favicon/favicon.png")
}
