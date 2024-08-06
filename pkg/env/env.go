package env

import (
	"log"
	"os"
)

//	get the environment var for the port you'd prefer to use via ENV in k8s/dockerfile,
//	or alternatively it will default to 80, if no port was specified

// EnvPort get the port from environment variables or return the default (80)
func EnvPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}

	return ":" + port
}

// EnvHostname get the hostname from environment variables or return the default (the hostname of the server)
func EnvHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err.Error())
		return os.Getenv("HOSTNAME")
	} else {
		return hostname
	}
}

// EnvCacheRefresh get the cache refresh from environment variables or return the default (60)
func EnvCacheRefresh() (string, int) {
	var score int

	timer := os.Getenv("REFRESHCACHE")

	if len(timer) == 0 {
		timer = "60"
		score = 1 // default
	} else {
		score = 2 // configured (doesn't matter what it is...)
	}

	return timer, score
}

// EnvTemplateDir get the templates directory from environment variables or return the default (/var/frontend/templates/)
func EnvTemplateDir() (string, int) {
	var score int

	templateDir := os.Getenv("TEMPLATEDIRECTORY")

	if len(templateDir) == 0 {
		templateDir = "/var/frontend/templates/"
		score = 1 // default
	} else {
		score = 2 // configured (doesn't matter what it is...)
	}

	return templateDir, score
}

// EnvTemplateType get the templates type from environment variables or return the default (.gohtml)
func EnvTemplateType() (string, int) {
	var score int

	templateType := os.Getenv("TEMPLATETYPE")

	if len(templateType) == 0 {
		templateType = "gohtml"
		score = 1
	} else {
		score = 2
	}

	return templateType, score
}

// EnvGitLink get the git repo from environment variables or return the default which is empty
func EnvGitLink() (string, int) {
	var score int

	gitlink := os.Getenv("GITLINK")

	if len(gitlink) == 0 {
		gitlink = ""
		score = 1
	} else {
		score = 2
	}

	return gitlink, score
}

// EnvGitToken get the git token from environment variables or return the default which is empty
func EnvGitToken() (string, int) {
	var score int

	gitToken := os.Getenv("GITTOKEN")

	if len(gitToken) == 0 {
		gitToken = ""
		score = 1
	} else {
		score = 2
	}

	return gitToken, score
}

// EnvGitUsername get the git username from environment variables or return the default which is empty
func EnvGitUsername() (string, int) {
	var score int

	gitUsername := os.Getenv("GITUSERNAME")

	if len(gitUsername) == 0 {
		gitUsername = ""
		score = 1
	} else {
		score = 2
	}

	return gitUsername, score
}

// EnvPathDefiner get the static paths from environment variables or return the default values based on switch expression
func EnvPathDefiner(path string) (string, int) {
	var score int
	var pathReturn string

	switch path {
	case "css":
		pathReturn = os.Getenv("CSSPATH")

		if len(pathReturn) == 0 {
			pathReturn = "css"
			score = 1
		} else {
			score = 2
		}

	case "js":
		pathReturn = os.Getenv("JSPATH")

		if len(pathReturn) == 0 {
			pathReturn = "js"
			score = 1
		} else {
			score = 2
		}

	case "images":
		pathReturn = os.Getenv("IMAGESPATH")

		if len(pathReturn) == 0 {
			pathReturn = "images"
			score = 1
		} else {
			score = 2
		}

	case "downloads":
		pathReturn = os.Getenv("DOWNLOADSPATH")

		if len(pathReturn) == 0 {
			pathReturn = "downloads"
			score = 1
		} else {
			score = 2
		}
	}

	return pathReturn, score
}
