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

	"poniatowski-dev/internal"
	"poniatowski-dev/internal/backendconnector"
	"poniatowski-dev/internal/caching"
	"poniatowski-dev/internal/env"
	"poniatowski-dev/internal/kubernetes"
	"poniatowski-dev/internal/logging"

	"github.com/joho/godotenv"
)

func main() {
	isReady := &atomic.Value{}
	isReady.Store(false)
	// Starting liveness probe
	fmt.Printf(" Starting liveness probe...")
	http.HandleFunc("/k8s/api/health", kubernetes.Healthz)
	fmt.Println(" Done.")
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
			timer, err := strconv.Atoi(env.EnvCacheRefresh())
			if err != nil {
				logging.LogIt("EnvCacheRefresh", "ERROR", "error loading environment variables")
			}
			time.Sleep(time.Duration(timer) * time.Second)
		}
	}()
	fmt.Println(" Done.")
	// JS and CSS handling/serving
	fmt.Printf(" Serving Static CSS/JS...")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(internal.TemplatesDir+"css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(internal.TemplatesDir+"js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(internal.TemplatesDir+"images"))))
	fmt.Println(" Done.")
	// Initialize the database and determine if collector should be enabled
	fmt.Printf(" Initializing Database...")
	errEnv := godotenv.Load("/usr/local/env/.env")
	if errEnv != nil {
		logging.LogIt("DatabaseExists", "ERROR", "error loading .env variables")
	}
	internal.CollectionError = backendconnector.DatabaseInit()
	if internal.CollectionError != nil {
    log.Println(" Failed.")
		logging.LogIt("Main", "WARNING", "Failed to initialize database, data will not be collected.")
	} else {
		fmt.Println(" Done.")
	}
	// Page serving section
	fmt.Printf(" Initializing Routes...")
	http.HandleFunc("/favicon.ico", internal.FaviconHandler) // favicon
	http.HandleFunc("/", cache.HandlerFunc)
	fmt.Println(" Done.")
	// Starting readiness probe
	fmt.Printf(" Starting readiness probe...")
	http.HandleFunc("/k8s/api/ready", kubernetes.Readyz(isReady))
	fmt.Println(" Done.")
	fmt.Printf(" Starting HTTP server on %v port %v...\n", env.EnvHostname(), strings.TrimLeft(env.EnvPort(), ":"))
	fmt.Println("<======================================================================================>")
	err := http.ListenAndServe(env.EnvPort(), nil)
	if err != nil {
		fmt.Println(" Error starting HTTP server...")
		fmt.Println("<======================================================================================>")
		log.Println(err.Error())
	}
}
