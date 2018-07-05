package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

// ClustersService performs operations on clusters.
type ClustersService interface {
	List(args ListArguments) (clusters ClustersResult, err error)
	Create(name string) (result Cluster, err error)
	Get(uuid string) (result Cluster, err error)
}

// GenericClustersService is a ClusterService placeholder implementation.
type GenericClustersService struct {
	connectionUrl string
}

// ListArguments are arguments relevant for listing objects.
type ListArguments struct {
	Page int
	Size int
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
	Name string
	UUID string
}

// NewClustersService Creates a new ClustersService.
func NewClustersService(connectionUrl string) ClustersService {
	service := new(GenericClustersService)
	service.connectionUrl = connectionUrl
	return service
}

// List returns lists of clusters.
func (cs GenericClustersService) List(args ListArguments) (result ClustersResult, err error) {
	result.Items = make([]Cluster, 0)
	db, err := sql.Open("postgres", cs.connectionUrl)
	if err != nil {
		return ClustersResult{}, fmt.Errorf("Error openning connection: %v", err)
	}
	defer db.Close()
	fmt.Printf("LIMIT: [%d] OFFESET: [%d]\n", args.Size, args.Page*args.Size)
	rows, err := db.Query(`SELECT uuid, name
		FROM clusters
		ORDER BY uuid
		LIMIT $1
		OFFSET $2`,
		args.Size,
		args.Page*args.Size,
	)
	if err != nil {
		return ClustersResult{}, fmt.Errorf("Error executing query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var uuid, name string
		err = rows.Scan(&uuid, &name)
		if err != nil {
			return ClustersResult{}, err
		}
		result.Items = append(result.Items, Cluster{
			Name: name,
			UUID: uuid,
		})
	}
	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		return ClustersResult{}, err
	}
	result.Page = args.Page
	result.Size = len(result.Items)
	return result, nil
}

// Create saves a new cluster definition in the Database
func (cs GenericClustersService) Create(name string) (result Cluster, err error) {
	uuid, err := ksuid.NewRandom()
	fmt.Printf("Generated id: %s\n", uuid)
	if err != nil {
		return Cluster{}, err
	}
	db, err := sql.Open("postgres", cs.connectionUrl)
	if err != nil {
		return Cluster{}, err
	}
	defer db.Close()
	stmt, err := db.Prepare(`INSERT INTO clusters (uuid, name)
		VALUES ($1, $2)`)
	if err != nil {
		return Cluster{}, err
	}
	defer stmt.Close()
	queryResult, err := stmt.Exec(uuid, name)
	if err != nil {
		return Cluster{}, err
	}
	inserted, err := queryResult.RowsAffected()
	if err != nil {
		return Cluster{}, err
	}
	if inserted != 1 {
		return Cluster{}, fmt.Errorf("Error: [%d] rows inserted. 1 expected",
			inserted,
		)
	}
	return Cluster{
		Name: name,
		UUID: fmt.Sprintf("%s", uuid),
	}, nil
}

// Get returns a single cluster by id
func (cs GenericClustersService) Get(uuid string) (result Cluster, err error) {
	db, err := sql.Open("postgres", cs.connectionUrl)
	if err != nil {
		return Cluster{}, err
	}
	defer db.Close()
	var name string
	err = db.QueryRow("SELECT uuid, name FROM clusters	WHERE uuid = $1", uuid).Scan(&uuid, &name)
	if err != nil {
		return Cluster{}, err
	}
	return Cluster{uuid, name}, nil
}
