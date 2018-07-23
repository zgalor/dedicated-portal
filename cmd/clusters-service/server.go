package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/container-mgmt/dedicated-portal/pkg/signals"
	"github.com/container-mgmt/dedicated-portal/pkg/sql"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var serveArgs struct {
	dbURL string
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the clusters service service",
	Long:  "Serve the clusters service service.",
	Run:   runServe,
}

// Server serves HTTP API requests on clusters.
type Server struct {
	stopCh         <-chan struct{}
	clusterService ClustersService
}

func init() {
	flags := serveCmd.Flags()
	flags.StringVar(
		&serveArgs.dbURL,
		"db-url",
		"",
		"The database connection url.",
	)
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
	apiRouter.HandleFunc("/clusters", s.listClusters).Methods(http.MethodGet)
	apiRouter.HandleFunc("/clusters", s.createCluster).Methods(http.MethodPost)
	apiRouter.HandleFunc("/clusters/{uuid}", s.getCluster).Methods(http.MethodGet)

	// Enable the access log:
	loggedRouter := handlers.LoggingHandler(os.Stdout, mainRouter)

	fmt.Println("Listening.")
	go http.ListenAndServe(":8000", loggedRouter)
	return nil
}

func runServe(*cobra.Command, []string) {
	// Set up signals so we handle the first shutdown signal gracefully:
	stopCh := signals.SetupHandler()

	// Check for db url cli arg:
	if serveArgs.dbURL == "" {
		glog.Errorf("flag missing: --db-url")
		os.Exit(1)
	}

	err := sql.EnsureSchema(
		"/usr/local/share/clusters-service/migrations",
		serveArgs.dbURL,
	)
	if err != nil {
		glog.Errorf("can't run sql migration: %s", err)
		os.Exit(1)
	}
	service := NewClustersService(serveArgs.dbURL)
	fmt.Println("Created cluster service.")

	// This is temporary and should be replaced with reading from the queue
	server := NewServer(stopCh, service)
	err = server.start()
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
	fmt.Println("Created server.")

	fmt.Println("Waiting for stop signal")
	<-stopCh // wait until requested to stop.
}
