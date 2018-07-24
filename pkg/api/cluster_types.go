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

// This file contains the API types used by the clusters service.

package api

// Cluster represents a cluster.
type Cluster struct {
	ID      string          `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Region  string          `json:"region,omitempty"`
	Nodes   ClusterNodes    `json:"nodes,omitempty"`
	Memory  ClusterResource `json:"memory,omitempty"`
	CPU     ClusterResource `json:"cpu,omitempty"`
	Storage ClusterResource `json:"storage,omitempty"`
}

// ClusterNodes represents the node count inside a cluster.
type ClusterNodes struct {
	Total   int `json:"total,omitempty"`
	Master  int `json:"master,omitempty"`
	Infra   int `json:"infra,omitempty"`
	Compute int `json:"compute,omitempty"`
}

// ClusterList is a list of clusters.
type ClusterList struct {
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int        `json:"total"`
	Items []*Cluster `json:"items"`
}

// ClusterResource represents a resource availability in the cluster.
type ClusterResource struct {
	Used  int `json:"used,omitempty"`
	Total int `json:"total,omitempty"`
}
