/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// Server serves REST API requests on clusters.
type Server struct {
	service CustomersService
}

var serveArgs struct {
	host              string
	port              int
	sqlConnStr        string
	notificationTopic string
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the customers service service",
	Long:  "Serve the customers service service.",
	Run:   runServe,
}

func init() {
	flags := serveCmd.Flags()
	flags.StringVar(
		&serveArgs.host,
		"host",
		"0.0.0.0",
		"The IP address or host name of the server.",
	)
	flags.IntVar(
		&serveArgs.port,
		"port",
		8000,
		"The port number of the server.",
	)
	flags.StringVar(
		&serveArgs.sqlConnStr,
		"sql-connection-string",
		"host=localhost port=5432 user=postgres password=1234 dbname=customers sslmode=disable",
		"The connection string for connection to sql datastore.",
	)
	flags.StringVar(
		&serveArgs.notificationTopic,
		"notifications-topic",
		"customers.notifications",
		"The name of the topic listening to notifications, for example: customers.notifications",
	)
}

// InitServer is a constructor for the Server struct.
func initServer(service CustomersService) (server *Server) {
	server = new(Server)
	server.service = service
	return server
}

func runServe(cmd *cobra.Command, args []string) {
	service, err := NewSQLCustomersService(serveArgs.sqlConnStr)
	if err != nil {
		panic(fmt.Sprintf("Can't connect to etcd: %v", err))
	}
	defer service.Close()

	// Create server URL.
	serverAddress := fmt.Sprintf("%s:%d", serveArgs.host, serveArgs.port)

	// Inform user we are starting.
	glog.Infof("Starting customers-service server at %s.", serverAddress)

	// Start server.
	server := initServer(service)
	defer server.Close()

	// Create the main router:
	mainRouter := mux.NewRouter()

	// Create the API router:
	apiRouter := mainRouter.PathPrefix("/api/customers_mgmt/v1").Subrouter()
	apiRouter.HandleFunc("/customers", server.getCustomersList).Methods("GET")
	apiRouter.HandleFunc("/customers", server.addCustomer).Methods("POST")
	apiRouter.HandleFunc("/customers/{id}", server.getCustomerByID).Methods("GET")
	apiRouter.Path("/customers").
		Queries("page", "{[0-9]+}", "size", "{[0-9]+}").
		Methods("GET").
		HandlerFunc(server.getCustomersList)

	log.Fatal(http.ListenAndServe(serverAddress, mainRouter))
}

// Close server
func (server *Server) Close() {
	server.service.Close()
}
