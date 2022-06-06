package main

import (
	cs "ARS_Projekat/configstore"
	"ARS_Projekat/tracer"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeConfigBody(ctx context.Context, r io.Reader) (*cs.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeConfigBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var config *cs.Config
	if err := dec.Decode(&config); err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return config, nil
}

func decodeGroupBody(ctx context.Context, r io.Reader) (*cs.Group, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeGroupBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var group *cs.Group
	if err := dec.Decode(&group); err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return group, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}, id string) {
	span := tracer.StartSpanFromContext(ctx, "renderJSON")
	defer span.Finish()

	js, err := json.Marshal(v)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId(ctx context.Context) string {
	return uuid.New().String()
}
