package caching

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"poniatowski-dev/internal"
	"poniatowski-dev/internal/datacapturing"
	"poniatowski-dev/internal/logging"
	"strings"
)

// handlerFunc 1 function to handle all routes.
func (pages *Pages) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	//Collect request data
	if internal.CollectionError != nil {
		log.Println(internal.CollectionError)
	} else {
		go datacapturing.CollectData(r)
	}
	// Handle command line curl/wget and give their IP address back.
	if strings.Contains(r.Header.Get("User-Agent"), "curl") ||
		strings.Contains(r.Header.Get("User-Agent"), "Wget") ||
		strings.Contains(r.Header.Get("User-Agent"), "WindowsPowerShell") {
		ipAddress, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logging.LogIt("handlerFunc", "ERROR", "user_ip: "+r.RemoteAddr+" is not IP:port")
			_, ipErr := fmt.Fprintf(w, "userip: %q is not IP:port\n", r.RemoteAddr)
			if ipErr != nil {
				log.Println(ipErr)
				return
			}
		}
		userIP := net.ParseIP(ipAddress)
		if userIP == nil {
			_, err := fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
			if err != nil {
				log.Println(err)
			}
			logging.LogIt("handlerFunc", "ERROR", "user_ip: "+r.RemoteAddr+" is not IP:port")
			return
		}
		_, errWrite := fmt.Fprintf(w, "Your IP is: %s\n", ipAddress)
		if errWrite != nil {
			log.Println(errWrite)
		}
		forward := r.Header.Get("X-FORWARDED-FOR")
		if forward != "" {
			_, errWrite := fmt.Fprintf(w, "Forwarded for: %s\n", forward)
			if err != nil {
				log.Println(errWrite)
			}
		}
		return
	}
	w.Header().Set("Content-Type", "text/html")
  pageNotFound := false
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		logging.LogIt("405", "ERROR", "user_ip ["+r.RemoteAddr+"] : Method not allowed")
	} else {
		templates, ok := pages.RouteMap[strings.Trim(r.URL.Path, "/")]
		if r.URL.Path == "/" || r.URL.Path == "" {
			templates = pages.RouteMap["welcome"]
		} else {
			if !ok {
				if pages.RouteMap["404"] == nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					logging.LogIt("500", "ERROR", " --> welcome."+internal.TemplateType+" not found")
					return
				}
				templates = pages.RouteMap["404"]
				logging.LogIt("404", "ERROR", "user_ip ["+r.RemoteAddr+"] : Page requested not found")
        pageNotFound = true
			}
		}
    if pageNotFound {
      w.WriteHeader(http.StatusNotFound)
    }
		tmplErr := templates.ExecuteTemplate(w, "layout", nil)
		if tmplErr != nil {
			log.Println(tmplErr.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logging.LogIt("500", "ERROR", "user_ip ["+r.RemoteAddr+"] : Issue sending 'layout'")
		}
	}
}
