package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type Service struct {
	data map[string][]*Config //this is currently a database
}

func (ts *Service) createConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil || len(rt) == 1 {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	id := createId()
	ts.data[id] = rt
	renderJSON(w, ts.data)
}
func (ts *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil || len(rt) > 1 {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	id := createId()
	ts.data[id] = rt
	renderJSON(w, ts.data)
}

func (ts *Service) getAllConfigAndGroupHandler(w http.ResponseWriter, req *http.Request) {
	renderJSON(w, ts.data)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok || len(task) > 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, ts.data[id])
}

func (ts *Service) getConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok || len(task) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, ts.data[id])
}

func (ts *Service) deleteConfigAndGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *Service) putConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	_, ok := ts.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil || len(rt) < 1 {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	ts.data[id] = append(ts.data[id], rt...)
	renderJSON(w, ts.data[id])
}