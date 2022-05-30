package main

import (
	cs "ARS_Projekat/configstore"
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

	store, err := cs.New()

	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		store: store,
	}

	router.HandleFunc("/config/", countPostConfig(server.createConfigHandler)).Methods("POST")
	router.HandleFunc("/config/{id}/", countGetConfigVersion(server.getConfigVersionsHandler)).Methods("GET")
	router.HandleFunc("/config/{id}", countPostConfigVersion(server.putNewConfigVersion)).Methods("POST")
	router.HandleFunc("/config/{id}/{ver}/", countGetConfig(server.getConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/{ver}", countDeleteConfig(server.deleteConfigHandler)).Methods("DELETE")

	router.HandleFunc("/group/", countPostGroup(server.createGroupHandler)).Methods("POST")
	router.HandleFunc("/group/{id}", countPostGroupVersion(server.putNewGroupVersion)).Methods("POST")
	router.HandleFunc("/group/{id}/{ver}/", countGetGroup(server.getGroupHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/{ver}/", countDeleteGroup(server.deleteGroupHandler)).Methods("DELETE")
	router.HandleFunc("/group/{id}/{ver}/config/", countGetGroupConfigs(server.getConfigFromGroup)).Methods("GET")
	router.HandleFunc("/group/{id}/{ver}/config/", countAddGroupConfig(server.addConfigToGroupHandler)).Methods("POST")

	router.Path("/metrics").Handler(metricsHandler())

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
