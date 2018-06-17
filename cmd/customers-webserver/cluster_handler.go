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
)

// ClusterHandler returns an index of all clusters in the system
func ClusterHandler(w http.ResponseWriter, r *http.Request) {
	pageParam, ok := r.URL.Query()["page"]

	if !ok || len(pageParam) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing query parameter: page\n")
		return
	}
	page, err := strconv.Atoi(pageParam[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad query param: page, %v\n", err)
		return
	}

	sizeParam, ok := r.URL.Query()["size"]

	if !ok || len(sizeParam) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing query parameter: size\n")
		return
	}

	size, err := strconv.Atoi(sizeParam[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad query param: size, %v\n", err)
		return
	}

	ret, err := json.Marshal(MockGetClusters(page, size))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshal error: %v\n", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(ret)
	}
}
