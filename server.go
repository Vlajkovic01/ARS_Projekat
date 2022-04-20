package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type postServer struct {
	data map[string][]*Config
}

func (ts *postServer) createPostHandler() {}

func (ts *postServer) getAllHandler() {}

func (ts *postServer) getPostHandler() {}

func (ts *postServer) delPostHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
