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
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getQueryParamInt(key string, defaultValue int64, r *http.Request) (value int64, err error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue, nil
	}
	value, err = strconv.ParseInt(valStr, 10, 64)
	return
}

func (server *Server) getCustomersList(w http.ResponseWriter, r *http.Request) {
	var ret *CustomersList
	var err error
	var page int64
	var size int64

	// Get Query Parameters.
	page, err = getQueryParamInt("page", 0, r)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error listing customers, %v", err)})
		return
	}

	size, err = getQueryParamInt("size", defaultLimit, r)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error listing customers, %v", err)})
		return
	}

	args := &ListArguments{
		Page: page,
		Size: size,
	}

	ret, err = server.service.List(args)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error listing customers, %v", err)})
		return
	}
	writeJSONResponse(w, http.StatusOK, ret)

}

func (server *Server) addCustomer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var customer Customer
	err := decoder.Decode(&customer)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error decoding customer, %v", err)})
		return
	}
	ret, err := server.service.Add(customer)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error adding customer, %v", err)})
	} else {
		writeJSONResponse(w, http.StatusOK, ret)
	}
}

func (server *Server) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ret, err := server.service.Get(id)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error getting customer, %v", err)})
	} else {
		writeJSONResponse(w, http.StatusOK, ret)
	}
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		response, err = json.Marshal(map[string]string{"error": fmt.Sprint(err)})
		if err != nil {
			panic(fmt.Sprintf("error marshalling error: %v", err))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
