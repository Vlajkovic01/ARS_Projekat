package main

import (
	cs "ARS_Projekat/configstore"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeConfigBody(r io.Reader) (*cs.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var config *cs.Config
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func decodeGroupBody(r io.Reader) (*cs.Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var group *cs.Group
	if err := dec.Decode(&group); err != nil {
		return nil, err
	}
	return group, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//w.Write([]byte("Idempotence key: " + id))
}

func createId() string {
	return uuid.New().String()
}
