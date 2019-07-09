package sdk

// Server represents the server table in DB.
type Server struct {
	ID           string `json:"ID"`
	URL          string `json:"URL"`
	Name         string `json:"Name"`
	SerialNumber string `json:"SerialNumber"`
	Warnings     int    `json:"Warnings"`
	Criticals    int    `json:"Criticals"`
}

// ServerCollection represents a server collection.
type ServerCollection struct {
	Member []Server `json:"Member"`
}
