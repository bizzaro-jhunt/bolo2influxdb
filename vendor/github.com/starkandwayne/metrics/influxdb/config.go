package influxdb

type Config struct {
	Addr               string `json:"url"`
	User               string `json:"user"`
	Password           string `json:"password"`
	InsecureSkipVerify bool
	Database           string `json:"database"`
}
