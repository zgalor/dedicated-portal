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

	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

// Server serves REST API requests on clusters.
type Server struct {
	service CustomersService
	router  *mux.Router
	conn    client.Connection
}

var serveArgs struct {
	host              string
	port              int
	etcdEndpoint      string
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
		&serveArgs.etcdEndpoint,
		"etcd-endpoint",
		"localhost:2379",
		"The endpoint running the etcd data store.",
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
	server.router = mux.NewRouter()
	server.service = service
	return server
}

func runServe(cmd *cobra.Command, args []string) {
	service, err := NewEtcdCustomersService(serveArgs.etcdEndpoint)
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
	initConnection(server)
	defer server.Close()

	// Init route table
	server.router.HandleFunc("/customers", server.getAllCustomers).Methods("GET")
	server.router.HandleFunc("/customers", server.addCustomer).Methods("POST")
	server.router.HandleFunc("/customers/{id}", server.getCustomerByID).Methods("GET")

	log.Fatal(http.ListenAndServe(serverAddress, server.router))
}

func initConnection(server *Server) {
	// Create a new connection object.
	var err error

	server.conn, err = stomp.NewConnection(&client.ConnectionSpec{})
	if err != nil {
		log.Fatal(err)
	}
}

// Close server
func (server *Server) Close() {
	server.conn.Close()
	server.service.Close()
}
