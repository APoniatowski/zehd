package backendconnector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"poniatowski-dev/internal/logging"
)

type DatabaseExistsInfo struct {
	Frontend   string `json:"frontend"`
	Connection string `json:"connection"`
	Tables     string `json:"tables"`
}

type RequestData struct {
	FrontendName string `json:"frontendName"`
	TimeDate     int64  `json:"timeDate"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Path         string `json:"path"`
	Method       string `json:"method"`
	XForwardFor  string `json:"XForwardFor"`
	XRealIP      string `json:"XRealIP"`
	UserAgent    string `json:"useragent"`
	Via          string `json:"via"`
	Age          string `json:"age"`
	CFIPCountry  string `json:"CF-IPCountry"`
}

// DatabaseInit - Initialize the database with tables, if the check returns false
func DatabaseInit() error {
	var databaseInfo DatabaseExistsInfo
	databaseInfo.DatabaseExists()
	if databaseInfo.Tables == "exists" {
		logging.LogIt("DatabaseInit", "INFO", "database already exists")
	} else {
		switch databaseInfo.Connection {
		case "successful":
			logging.LogIt("DatabaseInit", "INFO", "database does not exist, creating database and tables")
			databaseInfo.DatabaseCreate()
			if databaseInfo.Tables == "created" {
				logging.LogIt("DatabaseInit", "INFO", "database and tables have been created")
			} else {
				logging.LogIt("DatabaseInit", "ERROR", "failed to create Database: "+databaseInfo.Tables)
			}

		case "failed":
			logging.LogIt("DatabaseInit", "ERROR", "unable to reach backend")

		case "no env var":
			logging.LogIt("DatabaseInit", "WARNING", "no backend to send data to, please add backend url as env variable")

		default:
			logging.LogIt("DatabaseInit", "ERROR", "an error occurred during db checks: "+databaseInfo.Connection)
			return errors.New("an error occurred during db checks: " + databaseInfo.Connection)
		}
	}
	return nil
}

// DatabaseExists - Check if database exists and has existing tables
func (dbInfo *DatabaseExistsInfo) DatabaseExists() {
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		dbInfo.Connection = "no env var"
		return
	}
	response, errHTTPGet := http.Get(backendURL + "/database/exist")
	if errHTTPGet != nil {
		dbInfo.Connection = "failed"
		logging.LogIt("DatabaseExists", "ERROR", "unable to reach backend on "+backendURL)
		return
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			errCloseString := fmt.Sprintf("%v", errClose)
			dbInfo.Connection = errCloseString
			logging.LogIt("DatabaseExists", "ERROR", "unable to close response body")

		}
	}()
	body, _ := io.ReadAll(response.Body)
	dbInfo.Tables = string(body) // TODO reading results from here
	logging.LogIt("DatabaseInit", "INFO", "received: \""+dbInfo.Tables+"\" from "+backendURL+"(GET request)")
}

func (dbInfo *DatabaseExistsInfo) DatabaseCreate() {
	// create post request
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		dbInfo.Connection = "no env var"
		logging.LogIt("DatabaseCreate", "WARNING", "no backend to send data to, please add backend url as env variable")
		return
	}
	dbInfo.Tables = "create"
	req, _ := json.Marshal(dbInfo)
	response, errHTTPGet := http.Post(backendURL+"/database/exist", "application/json", bytes.NewBuffer(req))
	if errHTTPGet != nil {
		logging.LogIt("DatabaseCreate", "ERROR", "unable to reach backend")
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			logging.LogIt("DatabaseCreate", "ERROR", "unable to close response body")
		}
	}()
}

// DBConnector  Function to insert request data into the database
func (rD *RequestData) DBConnector(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	jsonToBackend, errMarshal := json.Marshal(rD)
	if errMarshal != nil {
		logging.LogIt("DBConnector", "ERROR", "unable to marshal json request")
	}
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		logging.LogIt("DBConnector", "WARNING", "no backend to send data to, please add backend url as env variable")
	} else {
		resp, errResp := http.Post(backendURL+"/api/collect", "application/json", bytes.NewBuffer(jsonToBackend))
		if errResp != nil {
			logging.LogIt("DBConnector", "ERROR", "unable to send json request, response error received")
		}
		defer func() {
			errClose := resp.Body.Close()
			if errClose != nil {
				logging.LogIt("DBConnector", "ERROR", "unable to close response body")
			}
		}()
		if resp.StatusCode == 200 {
			log.Println("User data sent to database successfully")
		} else if resp.StatusCode == 500 {
			logging.LogIt("DBConnector", "WARNING", "no backend to send data to, please add backend url as env variable")
		}
	}
	waitGroup.Done()
}
