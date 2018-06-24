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
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
	"strings"
)

// SQLCustomersService is a struct implementing the customer service interface,
// backed by an SQL database.
type SQLCustomersService struct {
	db *sql.DB
}

const defaultLimit = 1000

// NewSQLCustomersService is a constructor for the SQLCustomersService struct.
func NewSQLCustomersService(connStr string) (*SQLCustomersService, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	service := new(SQLCustomersService)
	service.db = db
	return service, nil
}

// Close closes the etcd customers service client.
func (service *SQLCustomersService) Close() {
	service.db.Close()
}

// Add adds a single customer to psql database.
func (service *SQLCustomersService) Add(customer Customer) (*Customer, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return nil, err
	}

	result := Customer{
		ID:   id.String(),
		Name: customer.Name,
	}

	if customer.OwnedClusters == nil {
		result.OwnedClusters = make([]string, 0)
	} else {
		result.OwnedClusters = customer.OwnedClusters
	}

	_, err = service.db.Exec(`
		insert into customers (
		id,
		name
		) values (
			$1,
			$2
		)`,
		result.ID,
		result.Name)
	if err != nil {
		return nil, err
	}

	for _, cluster := range result.OwnedClusters {
		_, err = service.db.Exec(`
			insert into owned_clusters (
				customer_id,
				cluster_id
			) values (
				$1,
				$2
			)`,
			result.ID,
			cluster)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

// Get retrieves a single customer from psql database.
func (service *SQLCustomersService) Get(id string) (*Customer, error) {
	var result Customer

	// Get the customer information
	// If not customer found return nil pointer and nil error.
	// (See customers_service.go for more details)
	rows, err := service.db.Query(`select name from customers where id=$1`, id)
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		if err = rows.Scan(&result.Name); err != nil {
			return nil, err
		}
	}
	rows.Close()

	// Retrieve customer owned clusters.
	ownedClusters := make([]string, 0)
	rows, err = service.db.Query(`select cluster_id from owned_clusters
		where customer_id=$1`,
		id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var clusterID string
		if err = rows.Scan(&clusterID); err != nil {
			return nil, err
		}
		ownedClusters = append(ownedClusters, clusterID)
	}
	rows.Close()

	result.OwnedClusters = ownedClusters
	return &result, nil
}

// List retrieves a list of current customers stored in datastore.
func (service *SQLCustomersService) List(args *ListArguments) (*CustomersList, error) {
	var result *CustomersList
	var rows *sql.Rows
	var err error
	var page int64
	var numOfItems int64

	if args != nil {
		page = args.Page
		numOfItems = args.Size
	} else {
		page = 0
		numOfItems = defaultLimit
	}

	// Retrieve customers id's and names.
	rows, err = service.db.Query(`select id, name from customers
		 limit $1 offset $2`,
		numOfItems, numOfItems*page)
	if err != nil {
		return nil, err
	}

	// Populate customers id's and names in their corresponding customers struct.
	items := make([]*Customer, 0, numOfItems)
	ids := make([]string, 0, numOfItems)
	for rows.Next() {
		var customer Customer
		if err = rows.Scan(&customer.ID, &customer.Name); err != nil {
			return nil, err
		}
		// Populate items with customer id and customer names.
		items = append(items, &customer)
		// Keep id's to query for owned_clusters.
		ids = append(ids, customer.ID)
	}
	rows.Close()

	qoutedIds := make([]string, len(ids))
	for i, id := range ids {
		qoutedIds[i] = fmt.Sprintf(`'%s'`, id)
	}
	idSet := strings.Join(qoutedIds, ", ")
	query := fmt.Sprintf(`
		select customer_id, cluster_id
		from owned_clusters
		where customer_id in (%s)`,
		idSet)

	// Retrieve customers owned clusters.
	customersToClusters := make(map[string][]string)
	rows, err = service.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var clusterID string
		var customerID string
		if err = rows.Scan(&customerID, &clusterID); err != nil {
			return nil, err
		}
		customersToClusters[customerID] = append(customersToClusters[customerID], clusterID)
	}
	rows.Close()

	// Populate customers owned clusters
	for _, customer := range items {
		if customer != nil {
			customer.OwnedClusters = (customersToClusters[customer.ID])
		}
	}

	total, err := service.getCustomersCount()
	if err != nil {
		return nil, err
	}

	result = &CustomersList{
		Items: items,
		Page:  page,
		Size:  int64(len(items)),
		Total: total,
	}

	return result, nil
}

func (service *SQLCustomersService) getCustomersCount() (int64, error) {
	// retrieve total number of customers.
	var total int64
	err := service.db.QueryRow("select  count(*) from customers").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
