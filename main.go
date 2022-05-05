package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := Service{
		data: map[string][]*Config{},
	}
	router.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/group/", server.createConfigGroupHandler).Methods("POST")
	router.HandleFunc("/configs/", server.getAllConfigsHandler).Methods("GET")
	router.HandleFunc("/groups/", server.getAllGroupsHandler).Methods("GET")
	router.HandleFunc("/config/{id}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/config/{idConfig}/", server.getConfigFromGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.putConfigHandler).Methods("PUT")
	router.HandleFunc("/config/{id}/", server.deleteConfigHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/", server.deleteGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/config/{idConfig}/", server.deleteConfigFromGroupHandler).Methods("DELETE")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
