package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"zehd/pkg"
	"zehd/pkg/backendconnector"
	"zehd/pkg/caching"
	"zehd/pkg/env"
	"zehd/pkg/kubernetes"

	"github.com/APoniatowski/boillog"

	"github.com/joho/godotenv"
)

func main() {
	isReady := &atomic.Value{}
	isReady.Store(false)

	// Starting liveness probe
	fmt.Printf(" Starting liveness probe...")
	http.HandleFunc("/k8s/api/health", kubernetes.Healthz)
	fmt.Println(" Done.")

	// cloning repo from git source, if specified.
	fmt.Println(" Checking git repo:")
	fmt.Printf("  Cloning")

	if len(pkg.GitLink) == 0 {
		fmt.Printf("... Skipped.\n")
	} else {
		err := caching.Git("startup")
		if err != nil {
			fmt.Printf("\n")
			log.Println(err.Error())
		} else {
			fmt.Println(" Done.")
		}
	}

	// init caching via go func(){check files for changes and figure out a way to pass them to handler}
	fmt.Printf(" Building and Caching templates...")
	cache := &caching.Pages{}
	cache.RouteMap = make(map[string]*template.Template)
	go func() {
		for {
			err := cache.CachePages()
			if err != nil {
				log.Println(err.Error())
			}
			timer, err := strconv.Atoi(pkg.RefreshCache)
			if err != nil {
				boillog.LogIt("EnvCacheRefresh", "ERROR", "error loading environment variables")
			}
			time.Sleep(time.Duration(timer) * time.Second)
		}
	}()
	fmt.Println(" Done.")

	// JS and CSS handling/serving
	fmt.Printf(" Serving Static paths:\n")
	fmt.Printf("   CSS...")
	if pkg.CSS != pkg.Disable {
		http.Handle("/"+pkg.CSS+"/", http.StripPrefix("/"+pkg.CSS+"/", http.FileServer(http.Dir(pkg.TemplatesDir+pkg.CSS))))
		fmt.Println(" Done.")
	} else {
		fmt.Println(" Disabled.")
	}

	fmt.Printf("   JS... ")
	if pkg.Javascript != pkg.Disable {
		http.Handle("/"+pkg.Javascript+"/", http.StripPrefix("/"+pkg.Javascript+"/", http.FileServer(http.Dir(pkg.TemplatesDir+pkg.Javascript))))
		fmt.Println(" Done.")
	} else {
		fmt.Println(" Disabled.")
	}

	fmt.Printf("   Images... ")
	if pkg.Images != pkg.Disable {
		http.Handle("/"+pkg.Images+"/", http.StripPrefix("/"+pkg.Images+"/", http.FileServer(http.Dir(pkg.TemplatesDir+pkg.Images))))
		fmt.Println(" Done.")
	} else {
		fmt.Println(" Disabled.")
	}

	fmt.Printf("   Downloads... ")
	if pkg.Downloads != pkg.Disable {
		http.Handle("/"+pkg.Downloads+"/", http.StripPrefix("/"+pkg.Downloads+"/", http.FileServer(http.Dir(pkg.TemplatesDir+pkg.Downloads))))
		fmt.Println(" Done.")
	} else {
		fmt.Println(" Disabled.")
	}

	// Initialize the database and determine if collector should be enabled
	fmt.Printf(" Initializing Database...")
	errEnv := godotenv.Load("/usr/local/env/.env")
	if errEnv != nil {
		boillog.LogIt("DatabaseExists", "ERROR", "error loading .env variables")
	}

	pkg.CollectionError = backendconnector.DatabaseInit()
	if pkg.CollectionError != nil {
		log.Println(" Failed.")
		boillog.LogIt("Main", "WARNING", "Failed to initialize database, data will not be collected.")
	} else {
		fmt.Println(" Done.")
	}

	// Page serving section
	fmt.Printf(" Initializing Routes...")
	http.HandleFunc("/favicon.ico", pkg.FaviconHandler) // favicon
	http.HandleFunc("/", cache.HandlerFunc)
	fmt.Println(" Done.")

	// Starting readiness probe
	fmt.Printf(" Starting readiness probe...")
	http.HandleFunc("/k8s/api/ready", kubernetes.Readyz(isReady))
	fmt.Println(" Done.")

	// Starting server
	fmt.Printf(" Starting HTTP server on %v port %v...\n", env.EnvHostname(), strings.TrimLeft(env.EnvPort(), ":"))
	fmt.Println("<======================================================================================>")
	err := http.ListenAndServe(env.EnvPort(), nil)
	if err != nil {
		fmt.Println(" Error starting HTTP server...")
		fmt.Println("<======================================================================================>")
		log.Println(err.Error())
	}
}
