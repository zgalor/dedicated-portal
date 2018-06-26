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
	"strconv"

	"github.com/gorilla/mux"
)

// This is similar to the function in customers-webserver/utils.go
// there's no point turning it into a shared file between the two workers
// if the REST implementation is supposed to be removed later anyway.

func getQueryParamInt(param string, defaultValue int64, r *http.Request) (value int64, err error) {
	valueString, ok := r.URL.Query()[param]

	if !ok || len(valueString) < 1 {
		return defaultValue, nil
	}

	// This needs to be ParseInt and not Atoi because the interface asks for int64
	value, err = strconv.ParseInt(valueString[0], 10, 64)
	return
}

func (server *Server) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /customers request")

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

	ret, err := server.service.List(&ListArguments{Page: page, Size: size})
	if err != nil {
		writeJSONResponse(w, 500, map[string]string{"error": fmt.Sprintf("Error listing customers, %v", err)})
	} else {
		writeJSONResponse(w, 200, ret)
	}
}

func (server *Server) addCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /customers request")
	decoder := json.NewDecoder(r.Body)
	var customer Customer
	err := decoder.Decode(&customer)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error decoding customer, %v", err)})
		return
	}
	ret, err := server.service.Add(customer)
	if err != nil {
		writeJSONResponse(w, 500, map[string]string{"error": fmt.Sprintf("Error adding customer, %v", err)})
	} else {
		writeJSONResponse(w, 200, ret)
	}
}

func (server *Server) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /customers/{id} request")
	id := mux.Vars(r)["id"]
	ret, err := server.service.Get(id)
	if err != nil {
		writeJSONResponse(w, 500, map[string]string{"error": fmt.Sprintf("Error getting customer, %v", err)})
	} else {
		writeJSONResponse(w, 200, ret)
	}
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
