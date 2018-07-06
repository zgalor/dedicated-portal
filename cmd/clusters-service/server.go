package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// This file should be removed and replaced with a queue.

// Server serves HTTP API requests on clusters.
type Server struct {
	stopCh         <-chan struct{}
	clusterService ClustersService
}

// NewServer creates a new server.
func NewServer(stopCh <-chan struct{}, clusterService ClustersService) *Server {
	server := new(Server)
	server.stopCh = stopCh
	server.clusterService = clusterService
	return server
}

func (s Server) start() error {
	// Create the main router:
	mainRouter := mux.NewRouter()

	// Create the API router:
	apiRouter := mainRouter.PathPrefix("/api/clusters_mgmt/v1").Subrouter()
	apiRouter.HandleFunc("/clusters", s.listClusters).Methods("GET")
	apiRouter.HandleFunc("/clusters", s.createCluster).Methods("POST")
	apiRouter.HandleFunc("/clusters/{uuid}", s.getCluster).Methods("GET")

	fmt.Println("Listening.")
	go http.ListenAndServe(":8000", mainRouter)
	return nil
}

func (s Server) listClusters(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	page, err := getQueryParamInt("page", 0, r)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	size, err := getQueryParamInt("size", 100, r)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	results, err := s.clusterService.List(ListArguments{Page: page, Size: size})
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	fmt.Fprintf(w, "%s\n", b)
}

func (s Server) createCluster(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	var spec Cluster
	err = json.Unmarshal(bytes, &spec)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	if spec.UUID != "" {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "uuid must be empty"})
		return
	}
	result, err := s.clusterService.Create(spec.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s\n", b)
}

func (s Server) getCluster(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "no uuid provided"})
		return
	}
	cluster, err := s.clusterService.Get(uuid)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	b, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	fmt.Fprintf(w, "%s\n", b)
}

func getQueryParamInt(param string, defaultValue int, r *http.Request) (value int, err error) {
	valueString, ok := r.URL.Query()[param]

	if !ok || len(valueString) < 1 {
		return defaultValue, nil
	}
	var result int64
	// This needs to be ParseInt and not Atoi because the interface asks for int64
	result, err = strconv.ParseInt(valueString[0], 10, 32)
	return int(result), err
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.MarshalIndent(payload, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func logRequest(r *http.Request) {
	fmt.Printf("[%s] [%s]\n", r.Method, r.URL)
}
