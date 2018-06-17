package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// This file should be removed and replaced with a queue.

// Server serves HTTP API requests on clusters.
type Server struct {
	stopCh         <-chan struct{}
	clusterService ClustersService
}

func NewServer(stopCh <-chan struct{}, clusterService ClustersService) *Server {
	server := new(Server)
	server.stopCh = stopCh
	server.clusterService = clusterService
	return server
}

func (s Server) start() error {
	r := mux.NewRouter()
	r.HandleFunc("/clusters", func(w http.ResponseWriter, r *http.Request) {
		results, err := s.clusterService.List(ListArguments{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", b)

	}).Methods("GET")
	fmt.Println("Listening.")
	go http.ListenAndServe(":8000", r)
	fmt.Println("Listened.")
	return nil
}
