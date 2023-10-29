package env

import (
	"log"
	"os"
)

//	get the environment var for the port you'd prefer to use via ENV in k8s/dockerfile,
//	or alternatively it will default to 80, if no port was specified

func EnvPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	return ":" + port
}

func EnvHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err.Error())
		return os.Getenv("HOSTNAME")
	} else {
		return hostname
	}
}

func EnvTemplateDir() string {
	templateDir := os.Getenv("TEMPLATEDIRECTORY")
	if len(templateDir) == 0 {
		templateDir = "/var/frontend/templates/"
	}
	return templateDir
}

func EnvTemplateType() string {
	templateType := os.Getenv("TEMPLATETYPE")
	if len(templateType) == 0 {
		templateType = "gohtml"
	}
	return templateType
}

func EnvCacheRefresh() string {
	timer := os.Getenv("REFRESHCACHE")
	if len(timer) == 0 {
		timer = "60"
	}
	return timer
}

func EnvProfiler() string {
	profiler := os.Getenv("PROFILER")
	if len(profiler) == 0 {
		profiler = "false"
	}
	return profiler
}
