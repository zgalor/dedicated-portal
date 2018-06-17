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

	"github.com/container-mgmt/dedicated-portal/pkg/signals"
)

func main() {
	// Set up signals so we handle the first shutdown signal gracefully:
	stopCh := signals.SetupHandler()

	service := NewClustersService()
	fmt.Println("Created cluster service.")

	// This is temporary and should be replaced with reading from the queue
	server := NewServer(stopCh, service)
	err := server.start()
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}

	fmt.Println("Created server")
}
