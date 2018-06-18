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
)

// Cluster represents a cluster
type Cluster struct {
	Name          string `json:"name"`
	ID            int    `json:"id"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Nodes         int    `json:"nodes"`
	// TODO - figure out the actual parameters a cluster should have
}

// MockGetClusters return a list of clusters.
// TODO write an actual GetClusters function to get them from the database
func MockGetClusters(page int, size int) []*Cluster {
	clusters := make([]*Cluster, size)
	for i := 0; i < size; i++ {
		clusterNumber := (page * size) + i
		clusters[i] = &Cluster{
			Name:          fmt.Sprintf("test%d", clusterNumber),
			CloudProvider: "AWS",
			Region:        "us-east-1",
			Nodes:         5,
			ID:            clusterNumber}
	}
	return clusters
}
