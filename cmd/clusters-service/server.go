package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/container-mgmt/dedicated-portal/pkg/auth"
	"github.com/container-mgmt/dedicated-portal/pkg/signals"
	"github.com/container-mgmt/dedicated-portal/pkg/sql"
)

var serveArgs struct {
	jwkCertURL string
	dbURL      string
	demoMode   bool
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the clusters service service",
	Long:  "Serve the clusters service service.",
	Run:   runServe,
}

var (
	clusterOperatorKubeAddress string
	clusterOperatorKubeConfig  string
)

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
	flags.StringVar(
		&serveArgs.jwkCertURL,
		"jwk-certs-url",
		"",
		"The url endpoint for the JWK certs.",
	)
	flags.BoolVar(
		&serveArgs.demoMode,
		"demo-mode",
		false,
		"Run in demo mode (no token needed).",
	)
	flags.StringVar(
		&clusterOperatorKubeConfig,
		"cluster-operator-kubeconfig",
		"",
		"Path to a Kubernetes client configuration file used to connect "+
			"to the cluster where the cluster operator is running. Only required when running "+
			"cluster-operator outside of the cluster where the clusters-service is running. .",
	)
	flags.StringVar(
		&clusterOperatorKubeAddress,
		"cluster-operator-master",
		"",
		"The address of the Kubernetes API server for the cluster where cluster operator is running."+
			"Overrides any value in the Kubernetes "+
			"configuration file. Only required when running cluster-operator outside of the cluster "+
			"where the clusters-service is running.",
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
	var loggedRouter http.Handler

	// Create the main router:
	mainRouter := mux.NewRouter()

	// Create the API router:
	apiRouter := mainRouter.PathPrefix("/api/clusters_mgmt/v1").Subrouter()
	apiRouter.HandleFunc("/clusters", s.listClusters).Methods(http.MethodGet)
	apiRouter.HandleFunc("/clusters", s.createCluster).Methods(http.MethodPost)
	apiRouter.HandleFunc("/clusters/{uuid}", s.getCluster).Methods(http.MethodGet)

	// If not in demo mode, check JWK and add a JWT middleware:
	//
	// When running on demo mode we want to bypass the JWT check
	// and serve mock data.
	if !serveArgs.demoMode {
		// Check for JWK cert cli arg:
		if serveArgs.jwkCertURL == "" {
			check(fmt.Errorf("flag missing: --jwk-certs-url"), "No cert URL defined")
		}

		// Enable the access authentication:
		authRouter, err := auth.Router(serveArgs.jwkCertURL, mainRouter)
		check(
			err,
			fmt.Sprintf(
				"Can't create authentication route using URL '%s'",
				serveArgs.jwkCertURL,
			),
		)

		// Enable the access log:
		loggedRouter = handlers.LoggingHandler(os.Stdout, authRouter)
	} else {
		// Enable the access log:
		loggedRouter = handlers.LoggingHandler(os.Stdout, mainRouter)
	}

	fmt.Println("Listening.")
	go http.ListenAndServe(":8000", loggedRouter)
	return nil
}

func runServe(cmd *cobra.Command, args []string) {
	// Set up signals so we handle the first shutdown signal gracefully:
	stopCh := signals.SetupHandler()

	// Check for db url cli arg:
	if serveArgs.dbURL == "" {
		glog.Errorf("flag missing: --db-url")
		os.Exit(1)
	}

	k8sConfig, err := retrieveKubeConfig()
	if err != nil {
		glog.Fatalf("An error occurred while trying to retrieve kubernetes configurations: %s", err)
	}

	err = sql.EnsureSchema(
		"/usr/local/share/clusters-service/migrations",
		serveArgs.dbURL,
	)
	if err != nil {
		glog.Errorf("can't run sql migration: %s", err)
		os.Exit(1)
	}

	provisioner, err := NewClusterOperatorProvisioner(k8sConfig)
	if err != nil {
		panic(fmt.Sprintf("Error starting clusters service: %v", err))
	}

	service := NewClustersService(serveArgs.dbURL, provisioner)
	fmt.Println("Created cluster service.")

	server := NewServer(stopCh, service)
	err = server.start()
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}

	fmt.Println("Created server.")

	fmt.Println("Waiting for stop signal")
	<-stopCh // wait until requested to stop.
}

func retrieveKubeConfig() (*rest.Config, error) {
	// Load the Kubernetes configuration:
	var k8sConfig *rest.Config

	kubeConfig, err := kubeConfigPath(clusterOperatorKubeConfig)
	if err == nil {
		// If error is nil, we have a valid kubeConfig file:
		k8sConfig, err = clientcmd.BuildConfigFromFlags(clusterOperatorKubeAddress, kubeConfig)
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading REST client configuration from file '%s': %s",
				kubeConfig, err,
			)
		}

		return k8sConfig, nil

	} else if kubeConfig == "" {
		// If kubeConfig is "", file is missing, in this case we will
		// try to use in-cluster configuration.
		glog.Info("Try to use the in-cluster configuration")
		k8sConfig, err = rest.InClusterConfig()
		// Catch in-cluster configuration error:
		if err != nil {
			return nil, fmt.Errorf("Error loading in-cluster REST client configuration: %s", err)
		}

		return k8sConfig, nil

	} else {
		// Catch all errors:
		return nil, fmt.Errorf("Error: %s", err)
	}
}

func kubeConfigPath(clusterOperatorKubeConfig string) (kubeConfig string, err error) {
	// The loading order follows these rules:
	// 1. If the –kubeconfig flag is set,
	// then only that file is loaded. The flag may only be set once.
	// 2. If $KUBECONFIG environment variable is set, use it.
	// 3. Otherwise, ${HOME}/.kube/config is used.
	var ok bool

	// Get the config file path
	if clusterOperatorKubeConfig != "" {
		kubeConfig = clusterOperatorKubeConfig
	} else {
		if kubeConfig, ok = os.LookupEnv("KUBECONFIG"); ok != true {
			kubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}
	}

	// Check config file:
	fInfo, err := os.Stat(kubeConfig)
	if os.IsNotExist(err) {
		// NOTE: If config file does not exist, assume using pod configuration.
		err = fmt.Errorf("The Kubernetes configuration file '%s' doesn't exist", kubeConfig)
		kubeConfig = ""
		return
	}

	// Check error codes.
	if fInfo.IsDir() {
		err = fmt.Errorf("The Kubernetes configuration path '%s' is a direcory", kubeConfig)
		return
	}
	if os.IsPermission(err) {
		err = fmt.Errorf("Can't open Kubernetes configuration file '%s'", kubeConfig)
		return
	}

	return
}

// Exit on error
func check(err error, msg string) {
	if err != nil {
		glog.Errorf("%s: %s", msg, err)
		os.Exit(1)
	}
}
