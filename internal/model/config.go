package model

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Zone2HR      int    `json:"zone2_hr,omitempty"`
}
