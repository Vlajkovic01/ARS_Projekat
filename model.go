package main

type Config struct {
	Id      string              `json:"id"`
	Entries []map[string]string `json:"entries"`
}
