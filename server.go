package main

import (
	cs "ARS_Projekat/configstore"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"net/url"
)

type Service struct {
	store *cs.ConfigStore
}

func (ts *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

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

	rt, err := decodeConfigBody(req.Body)
	if err != nil || rt.Version == "" || rt.Entries == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}

	config, err := ts.store.CreateConfig(rt)

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	w.Write([]byte("Config ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) putNewConfigVersion(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigBody(req.Body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	rt.ID = id
	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}

	config, err := ts.store.UpdateConfigVersion(rt)

	if err != nil {
		http.Error(w, "Given config version already exists! ", http.StatusBadRequest)
		return
	}

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	w.Write([]byte("Config ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindConfig(id, ver)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindConfVersions(id)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) createGroupHandler(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

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

	rt, err := decodeGroupBody(req.Body)
	if err != nil || rt.Version == "" || rt.Configs == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}

	group, err := ts.store.CreateGroup(rt)

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	w.Write([]byte("Group ID: " + group.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	task, ok := ts.store.FindGroup(id, ver)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) getConfigFromGroup(w http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	req.ParseForm()
	params := url.Values.Encode(req.Form)
	labels, err := ts.store.FindLabels(id, ver, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, labels)
}

func (ts *Service) putNewGroupVersion(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	requestId := req.Header.Get("x-idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeGroupBody(req.Body)
	if err != nil || rt.Version == "" || rt.Configs == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}

	rt.ID = id
	config, err := ts.store.UpdateGroupVersion(rt)

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId()
	}

	if err != nil {
		http.Error(w, "Given config version already exists! ", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Group ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) addConfigToGroupHandler(w http.ResponseWriter, r *http.Request) {

	requestId := r.Header.Get("x-idempotency-key")

	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]
	var configs []map[string]string
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := dec.Decode(&configs)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	configs, err = ts.store.AddLabelsToGroup(configs, id, ver)

	if err != nil {
		http.Error(w, "key not found", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(requestId) == true {

		http.Error(w, "Request has been already sent", http.StatusBadRequest)

		return
	}

	reqId := ts.store.SaveRequestId()

	renderJSON(w, configs)
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]
	_, err := ts.store.DeleteConfig(id, ver)
	if err != nil {
		http.Error(w, "Could not delete config", http.StatusBadRequest)
	}
}

func (ts *Service) deleteGroupHandler(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	ver := mux.Vars(request)["ver"]
	err := ts.store.DeleteGroup(id, ver)
	if err != nil {
		http.Error(writer, "Could not delete group", http.StatusBadRequest)
	}
}
