package backendconnector

// DatabaseExistsInfo Struct for checking if the DB exists
type DatabaseExistsInfo struct {
	Frontend   string `json:"frontend"`
	Connection string `json:"connection"`
	Tables     string `json:"tables"`
}

// RequestData Struct for checking and collecting user requests via headers
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
