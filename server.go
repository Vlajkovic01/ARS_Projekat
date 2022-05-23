package main

import (
	"ARS_Projekat/configstore"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type Service struct {
	data map[string][]*configstore.Config //this is currently a database
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
		http.Error(w, "Invalid JSON format. Must be more than 1 config.", http.StatusBadRequest)
		return
	}

	id := createId()
	if _, exists := ts.data[id]; exists {
		http.Error(w, "The same request has already been sent.", http.StatusBadRequest)
		return
	}
	ts.data[id] = rt
	renderJSON(w, id)
}

func (ts *Service) createConfigGroupNewVersionHandler(w http.ResponseWriter, req *http.Request) {
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
		http.Error(w, "Invalid JSON format. Must be more than 1 config.", http.StatusBadRequest)
		return
	}

	id := mux.Vars(req)["id"]
	renderJSON(w, id)
}

func (ts *Service) createConfigNewVersionHandler(w http.ResponseWriter, req *http.Request) {

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
		http.Error(w, "Invalid JSON format. Must be exactly 1 config.", http.StatusBadRequest)
		return
	}

	id := mux.Vars(req)["id"]
	task := ts.data[id]
	task = append(task, rt...)
	ts.data[id] = task
	renderJSON(w, id)
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
		http.Error(w, "Invalid JSON format. Must be exactly 1 config.", http.StatusBadRequest)
		return
	}

	id := createId()
	if _, exists := ts.data[id]; exists {
		http.Error(w, "The same request has already been sent.", http.StatusBadRequest)
		return
	}
	ts.data[id] = rt
	renderJSON(w, id)
}

func (ts *Service) getAllConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := make(map[string][]*configstore.Config)
	for k, v := range ts.data {
		if len(v) < 2 {
			allTasks[k] = v
		}
	}

	renderJSON(w, allTasks)
}

func (ts *Service) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := make(map[string][]*configstore.Config)
	for k, v := range ts.data {
		if len(v) >= 2 {
			allTasks[k] = v
		}
	}

	renderJSON(w, allTasks)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.data[id]

	if !ok || len(task) > 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if ok {
		for _, v := range task {
			if v.Version != version {
				err := errors.New("version not found")
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}
	}
	renderJSON(w, task)
}

func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.data[id]
	if !ok || len(task) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if ok {
		for _, v := range task {
			if v.Version != version {
				err := errors.New("version not found")
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}
	}
	renderJSON(w, task)
}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, ok := ts.data[id]

	if !ok || len(config) > 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if ok {
		for _, v := range config {
			if v.Version != version {
				err := errors.New("version not found")
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}
	}

	delete(ts.data, id)
	renderJSON(w, config)
}

func (ts *Service) deleteGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	group, ok := ts.data[id]

	if !ok || len(group) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		delete(ts.data, id)
		renderJSON(w, group)
	}
}

func (ts *Service) putConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	group, ok := ts.data[id]
	if !ok || len(group) < 2 {
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
