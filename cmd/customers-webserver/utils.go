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

	"github.com/golang/glog"
)

func writeErrorJSON(w http.ResponseWriter, httpErrorCode int, errorText string) {
	w.WriteHeader(httpErrorCode)
	ret, err := json.Marshal(map[string]string{"error": errorText})
	if err != nil {
		glog.Errorf("Can't marshal json for error response: %v", err)
		fmt.Fprintf(w, "{\"error\": \"Can't marshal json for error response\"}")
	} else {
		w.Write(ret)
	}

}

func getQueryParamString(param string, r *http.Request) (string, error) {
	valueString, ok := r.URL.Query()[param]

	if !ok || len(valueString) < 1 {
		return "", fmt.Errorf("Missing query parameter: %s", param)
	}
	return valueString[0], nil
}

func getQueryParamInt(param string, r *http.Request) (value int, err error) {
	valueString, err := getQueryParamString(param, r)
	if err != nil {
		return
	}

	value, err = strconv.Atoi(valueString)
	return
}
