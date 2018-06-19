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

	"github.com/golang/glog"
)

// ClusterHandler returns an index of all clusters in the system
func ClusterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page, err := getQueryParamInt("page", r)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("Bad query param: page, %v", err))
		return
	}

	size, err := getQueryParamInt("size", r)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("Bad query param: size, %v", err))
		return
	}

	ret, err := json.Marshal(map[string]interface{}{
		"page":  page,
		"size":  size,
		"total": 10000,
		"items": MockGetClusters(page, size)})

	if err != nil {
		glog.Errorf("Can't marshal json for cluster list response: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("Marshal error, %v", err))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(ret)
	}
}
