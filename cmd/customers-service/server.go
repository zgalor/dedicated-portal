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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/container-mgmt/dedicated-portal/cmd/customers-service/jwtcert"
)

// Server serves REST API requests on clusters.
type Server struct {
	service CustomersService
}

var serveArgs struct {
	host              string
	port              int
	jwkCertURL        string
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
		&serveArgs.jwkCertURL,
		"jwk-certs-url",
		"https://localhost/auth/certs",
		"The url endpoint for the JWK certs.",
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
	// Try to connect to SQLCustomersService
	service, err := NewSQLCustomersService(serveArgs.sqlConnStr)
	check(err, "Can't connect to sql service")
	defer service.Close()

	// Try to read the JWT public key object file.
	jwtCert, err := jwtcert.DownloadAsPEM(serveArgs.jwkCertURL)
	check(
		err,
		fmt.Sprintf(
			"Can't download JWK certificate from URL '%s'",
			serveArgs.jwkCertURL,
		),
	)

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

	// JWT Middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(jwtCert))
			return result, nil
		},
		ErrorHandler:  onAuthError,
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Enable the access authentication:
	authRouter := negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext))
	authRouter.UseHandler(mainRouter)

	// Enable the access log:
	loggedRouter := handlers.LoggingHandler(os.Stdout, authRouter)

	// ListenAndServe
	log.Fatal(http.ListenAndServe(serverAddress, loggedRouter))
}

// Close server
func (server *Server) Close() error {
	return server.service.Close()
}

// onAuthError returns an error json struct
func onAuthError(w http.ResponseWriter, r *http.Request, err string) {
	msg, _ := json.Marshal(map[string]string{"error": fmt.Sprint(err)})
	if msg == nil {
		msg = []byte("Unknown error while converting an error to json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	responseWriterWriteWithLog(w, msg)
}

// Panic on error
func check(err error, msg string) {
	if err != nil {
		glog.Errorf("%s: %s", msg, err)
		os.Exit(1)
	}
}
