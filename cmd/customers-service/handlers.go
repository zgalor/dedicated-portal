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
	"log"
	"net/http"
)

func (server *Server) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /customers request")
	writeJSONResponse(w, 200, nil)
}

func (server *Server) addCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /customers request")
	writeJSONResponse(w, 200, nil)
}

func (server *Server) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /customers/{id} request")
	writeJSONResponse(w, 200, nil)
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
