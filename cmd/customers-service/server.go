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
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Server serves REST API requests on clusters.
type Server struct {
	service CustomersService
	router  *mux.Router
}

// InitServer is a constructor for the Server struct.
func InitServer(service CustomersService) (server *Server) {
	server = new(Server)
	server.router = mux.NewRouter()
	server.service = service
	initRoutes(server)
	return server
}

func initRoutes(server *Server) {
	server.router.HandleFunc("/customers", server.getAllCustomers).Methods("GET")
	server.router.HandleFunc("/customers", server.addCustomer).Methods("POST")
	server.router.HandleFunc("/customers/{id}", server.getCustomerByID).Methods("GET")
}

// Run starts the service REST API
func (server *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", server.router))
}
