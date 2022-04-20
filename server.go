package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	data map[string][]*Config //this is currently a database
}

func (ts *Service) createConfigHandler() {}

func (ts *Service) getAllConfigHandler() {}

func (ts *Service) getConfigHandler() {}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
