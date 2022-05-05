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

	for _, v := range rt {
		v.Entries["id"] = createId()
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
	for _, v := range rt {
		v.Entries["id"] = id
	}

	ts.data[id] = append(ts.data[id], rt...)
	renderJSON(w, rt)
}

func (ts *Service) getAllConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		if len(v) < 2 {
			allTasks = append(allTasks, v...)
		}
	}

	renderJSON(w, allTasks)
}

func (ts *Service) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		if len(v) >= 2 {
			allTasks = append(allTasks, v...)
		}
	}

	renderJSON(w, allTasks)
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

func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok || len(task) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, ts.data[id])
}

func (ts *Service) getConfigFromGroupHandler(w http.ResponseWriter, req *http.Request) {
	idGroup := mux.Vars(req)["id"]
	idConfig := mux.Vars(req)["idConfig"]

	group, ok := ts.data[idGroup]

	if !ok || len(group) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, v := range group {
		if v.Entries["id"] == idConfig {
			renderJSON(w, v)
		}
	}
}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	config, ok := ts.data[id]
	if !ok || len(config) > 1 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		delete(ts.data, id)
		renderJSON(w, config)
	}
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

func (ts *Service) deleteConfigFromGroupHandler(w http.ResponseWriter, req *http.Request) {
	idGroup := mux.Vars(req)["id"]
	idConfig := mux.Vars(req)["idConfig"]

	group, ok := ts.data[idGroup]

	if !ok || len(group) < 2 {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		renderJSON(w, group)
		for _, v := range group {
			if v.Entries["id"] == idConfig {
				delete(ts.data, idConfig)
				renderJSON(w, v)
			}
		}
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

	for _, v := range rt {
		v.Entries["id"] = createId()
	}

	ts.data[id] = append(ts.data[id], rt...)
	renderJSON(w, ts.data[id])
}
