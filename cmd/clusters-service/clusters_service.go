package main

import "time"

// ClustersService performs operations on clusters.
type ClustersService interface {
	List(args ListArguments) (clusters ClustersResult, err error)
}

// GenericClustersService is a ClusterService placeholder implementation.
type GenericClustersService struct{}

// ListArguments are arguments relevant for listing objects.
type ListArguments struct {
	Skip  int
	Limit int
}

// ClustersResult is a result for a List request of Clusters.
type ClustersResult struct {
	Page  int
	Size  int
	Total int
	Items []Cluster
}

// Cluster represents an OpenShift cluster.
type Cluster struct {
	Name              string
	UUID              string
	CreationTimestamp time.Time
}

// NewClustersService Creates a new ClustersService.
func NewClustersService() ClustersService {
	return new(GenericClustersService)
}

// List returns lists of clusters.
func (cs GenericClustersService) List(args ListArguments) (result ClustersResult, err error) {
	result.Page = 0
	result.Size = 1
	result.Total = 1
	result.Items = []Cluster{
		Cluster{
			Name:              "Static Cluster",
			UUID:              "8a83de13-7a3e-4928-828b-74ffede43964",
			CreationTimestamp: time.Now(),
		},
	}
	return
}
