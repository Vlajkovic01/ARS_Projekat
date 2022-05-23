package configstore

type Group struct {
	ID      string              `json:"id"`
	Configs []map[string]string `json:"configs"`
	Version string              `json:"version"`
}

type Config struct {
	ID      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}
