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

// EnvTemplateDir get the templates directory from environment variables or return the default (/var/frontend/templates/)
func EnvTemplateDir() string {
	templateDir := os.Getenv("TEMPLATEDIRECTORY")
	if len(templateDir) == 0 {
		templateDir = "/var/frontend/templates/"
	}
	return templateDir
}

// EnvTemplateType get the templates type from environment variables or return the default (.gohtml)
func EnvTemplateType() string {
	templateType := os.Getenv("TEMPLATETYPE")
	if len(templateType) == 0 {
		templateType = "gohtml"
	}
	return templateType
}

// EnvCacheRefresh get the cache refresh from environment variables or return the default (60)
func EnvCacheRefresh() string {
	timer := os.Getenv("REFRESHCACHE")
	if len(timer) == 0 {
		timer = "60"
	}
	return timer
}
