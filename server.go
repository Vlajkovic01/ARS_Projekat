package main

import (
	cs "ARS_Projekat/configstore"
	"ARS_Projekat/tracer"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"io"
	"mime"
	"net/http"
	"net/url"
)

type Service struct {
	store  *cs.ConfigStore
	tracer opentracing.Tracer
	closer io.Closer
}

const (
	name = "configstore"
)

func NewConfigServer() (*Service, error) {
	store, err := cs.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(tracer)
	return &Service{
		store:  store,
		tracer: tracer,
		closer: closer,
	}, nil
}

func (s *Service) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *Service) GetCloser() io.Closer {
	return s.closer
}

func (s *Service) CloseTracer() error {
	return s.closer.Close()
}

func (ts *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling config create at %s\n", req.URL.Path)),
	)

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

	ctx := tracer.ContextWithSpan(context.Background(), span)

	rt, err := decodeConfigBody(ctx, req.Body)
	if err != nil || rt.Version == "" || rt.Entries == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusForbidden)
		return
	}

	config, err := ts.store.CreateConfig(ctx, rt)

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId(ctx)
	}
	w.Write([]byte("Config ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) putNewConfigVersion(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("putNewConfigVersion", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling create new config version at %s\n", req.URL.Path)),
	)

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

	ctx := tracer.ContextWithSpan(context.Background(), span)

	rt, err := decodeConfigBody(ctx, req.Body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	rt.ID = id

	if ts.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusForbidden)
		return
	}

	config, err := ts.store.UpdateConfigVersion(ctx, rt)

	if err != nil {
		http.Error(w, "Given config version already exists! ", http.StatusBadRequest)
		return
	}

	reqId := ""

	if err == nil {

		reqId = ts.store.SaveRequestId(ctx)
	}

	w.Write([]byte("Config ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling get config at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindConfig(ctx, id, ver)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(ctx, w, task, "")
}

func (ts *Service) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigVersionsHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling get config versions at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindConfVersions(ctx, id)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(ctx, w, task, "")
}

func (ts *Service) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createGroupHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling create group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

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

	rt, err := decodeGroupBody(ctx, req.Body)
	if err != nil || rt.Version == "" || rt.Configs == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusForbidden)
		return
	}

	group, err := ts.store.CreateGroup(ctx, rt)

	reqId := ""

	if err == nil {

		reqId = ts.store.SaveRequestId(ctx)
	}

	w.Write([]byte("Group ID: " + group.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling get group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	task, ok := ts.store.FindGroup(ctx, id, ver)

	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(ctx, w, task, "")
}

func (ts *Service) getConfigFromGroup(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigFromGroup", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling get config from group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	ver := mux.Vars(req)["ver"]
	id := mux.Vars(req)["id"]

	req.ParseForm()
	params := url.Values.Encode(req.Form)
	labels, err := ts.store.FindLabels(ctx, id, ver, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(ctx, w, labels, "")
}

func (ts *Service) putNewGroupVersion(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("putNewGroupVersion", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling put new group version at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

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

	rt, err := decodeGroupBody(ctx, req.Body)
	if err != nil || rt.Version == "" || rt.Configs == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusForbidden)
		return
	}

	rt.ID = id

	config, err := ts.store.UpdateGroupVersion(ctx, rt)

	reqId := ""

	if err == nil {
		reqId = ts.store.SaveRequestId(ctx)
	}

	if err != nil {
		http.Error(w, "Given config version already exists! ", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Group ID: " + config.ID))
	w.Write([]byte("\n\nIdempotence key: " + reqId))
}

func (ts *Service) addConfigToGroupHandler(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigToGroupHandler", ts.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling add config to group at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

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

	configs, err = ts.store.AddLabelsToGroup(ctx, configs, id, ver)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if ts.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusForbidden)
		return
	}

	reqId := ts.store.SaveRequestId(ctx)

	w.Write([]byte("Idempotence key: " + reqId))
}

func (ts *Service) deleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deleteConfigHandler", ts.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling delete config at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]
	_, err := ts.store.DeleteConfig(ctx, id, ver)
	if err != nil {
		http.Error(w, "Could not delete config", http.StatusBadRequest)
	}
}

func (ts *Service) deleteGroupHandler(writer http.ResponseWriter, request *http.Request) {
	span := tracer.StartSpanFromRequest("deleteGroupHandler", ts.tracer, request)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("Handling delete group at %s\n", request.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(request)["id"]
	ver := mux.Vars(request)["ver"]
	err := ts.store.DeleteGroup(ctx, id, ver)
	if err != nil {
		http.Error(writer, "Could not delete group", http.StatusBadRequest)
	}
}
