package main

type Service struct {
	Data map[string][]*Config `json:"config"`
}

type Config struct {
	Entries map[string]string `json:"entries"`
}
