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
	"time"

	"github.com/APoniatowski/boillog"
)

// DatabaseInit Initialize the database with tables, if the check returns false
func DatabaseInit() error {
	defer boillog.TrackTime("db-init", time.Now())
	var databaseInfo DatabaseExistsInfo
	databaseInfo.DatabaseExists()
	if databaseInfo.Tables == "exists" {
		boillog.LogIt("DatabaseInit", "INFO", "database already exists")
	} else {
		switch databaseInfo.Connection {
		case "successful":
			boillog.LogIt("DatabaseInit", "INFO", "database does not exist, creating database and tables")
			databaseInfo.DatabaseCreate()
			if databaseInfo.Tables == "created" {
				boillog.LogIt("DatabaseInit", "INFO", "database and tables have been created")
			} else {
				boillog.LogIt("DatabaseInit", "ERROR", "failed to create Database: "+databaseInfo.Tables)
			}

		case "failed":
			boillog.LogIt("DatabaseInit", "ERROR", "unable to reach backend")

		case "no env var":
			boillog.LogIt("DatabaseInit", "WARNING", "no backend to send data to, please add backend url as env variable")

		default:
			boillog.LogIt("DatabaseInit", "ERROR", "an error occurred during db checks: "+databaseInfo.Connection)
			return errors.New("an error occurred during db checks: " + databaseInfo.Connection)
		}
	}
	return nil
}

// DatabaseExists Check if database exists and has existing tables
func (dbInfo *DatabaseExistsInfo) DatabaseExists() {
	defer boillog.TrackTime("db-exists", time.Now())
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		dbInfo.Connection = "no env var"
		return
	}
	response, errHTTPGet := http.Get(backendURL + "/database/exist")
	if errHTTPGet != nil {
		dbInfo.Connection = "failed"
		boillog.LogIt("DatabaseExists", "ERROR", "unable to reach backend on "+backendURL)
		return
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			errCloseString := fmt.Sprintf("%v", errClose)
			dbInfo.Connection = errCloseString
			boillog.LogIt("DatabaseExists", "ERROR", "unable to close response body")

		}
	}()
	body, _ := io.ReadAll(response.Body)
	dbInfo.Tables = string(body) // TODO reading results from here
	boillog.LogIt("DatabaseInit", "INFO", "received: \""+dbInfo.Tables+"\" from "+backendURL+"(GET request)")
}

func (dbInfo *DatabaseExistsInfo) DatabaseCreate() {
	// create post request
	defer boillog.TrackTime("create-db", time.Now())
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		dbInfo.Connection = "no env var"
		boillog.LogIt("DatabaseCreate", "WARNING", "no backend to send data to, please add backend url as env variable")
		return
	}
	dbInfo.Tables = "create"
	req, _ := json.Marshal(dbInfo)
	response, errHTTPGet := http.Post(backendURL+"/database/exist", "application/json", bytes.NewBuffer(req))
	if errHTTPGet != nil {
		boillog.LogIt("DatabaseCreate", "ERROR", "unable to reach backend")
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			boillog.LogIt("DatabaseCreate", "ERROR", "unable to close response body")
		}
	}()
}

// DBConnector Function to insert request data into the database
func (rD *RequestData) DBConnector(waitGroup *sync.WaitGroup) {
	defer boillog.TrackTime("db-connector", time.Now())
	waitGroup.Add(1)
	jsonToBackend, errMarshal := json.Marshal(rD)
	if errMarshal != nil {
		boillog.LogIt("DBConnector", "ERROR", "unable to marshal json request")
	}
	backendURL := os.Getenv("BACKEND")
	if len(backendURL) == 0 {
		boillog.LogIt("DBConnector", "WARNING", "no backend to send data to, please add backend url as env variable")
	} else {
		resp, errResp := http.Post(backendURL+"/api/collect", "application/json", bytes.NewBuffer(jsonToBackend))
		if errResp != nil {
			boillog.LogIt("DBConnector", "ERROR", "unable to send json request, response error received")
		}
		defer func() {
			errClose := resp.Body.Close()
			if errClose != nil {
				boillog.LogIt("DBConnector", "ERROR", "unable to close response body")
			}
		}()
		if resp.StatusCode == 200 {
			log.Println("User data sent to database successfully")
		} else if resp.StatusCode == 500 {
			boillog.LogIt("DBConnector", "WARNING", "no backend to send data to, please add backend url as env variable")
		}
	}
	waitGroup.Done()
}
