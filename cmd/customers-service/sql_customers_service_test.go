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
	"os"
	"testing"
)

var service *SQLCustomersService
var localTesting = true

func TestMain(m *testing.M) {
	var err error
	if !localTesting {
		os.Exit(0)
	}
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=customers sslmode=disable"
	service, err = NewSQLCustomersService(connStr)
	if err != nil {
		fmt.Printf("An error occurred while trying to connect to database: %s\n", err)
	}
	defer deleteAll()
	defer service.Close()
	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	deleteAll()

	customerToAdd := Customer{
		Name: "test_customer",
	}
	res, err := service.Add(customerToAdd)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	if res.Name != "test_customer" {
		t.Fail()
	}
	if len(res.OwnedClusters) != 0 {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	deleteAll()

	var err error
	var customer *Customer

	expected := Customer{
		ID:            "customer-fake-id",
		Name:          "test_customer",
		OwnedClusters: []string{"fake-cluster-id0", "fake-cluster-id1"},
	}

	_, err = service.db.Exec("insert into customers (id, name) values ($1, $2)",
		expected.ID, expected.Name)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}

	for _, clusterID := range expected.OwnedClusters {
		_, err = service.db.Exec(`
			insert into owned_clusters
			(customer_id, cluster_id)
			values ($1, $2)`,
			expected.ID, clusterID)
		if err != nil {
			t.Fatal(err)
			t.Fail()
		}
	}

	customer, err = service.Get(expected.ID)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}

	if customer.Name != expected.Name {
		t.Fatalf("expected customer name to be: %s instead it was %s",
			expected.Name, customer.Name)
		t.Fail()
	}

	if len(customer.OwnedClusters) != len(expected.OwnedClusters) {
		t.Fatalf("expected customers number of clusters to be: %d instead it was %d",
			len(customer.OwnedClusters), len(expected.OwnedClusters))
		t.Fail()
	}
}

func TestList(t *testing.T) {
	deleteAll()

	items := []*Customer{
		{
			ID:            "test-id0",
			Name:          "test-name0",
			OwnedClusters: []string{"test-cluster-id0", "test-cluster-id1"},
		},
		{
			ID:            "test-id1",
			Name:          "test-name1",
			OwnedClusters: []string{"test-cluster-id2"},
		},
	}
	expected := CustomersList{
		Page:  1,
		Size:  2,
		Total: 2,
		Items: items,
	}

	for _, item := range expected.Items {
		_, err := service.db.Exec(`
			insert into customers
			 (id, name)
			 values ($1, $2)`,
			item.ID, item.Name)
		if err != nil {
			t.Fatal(err)
			t.Fail()
		}
		// retrieve customers owned clusters
		for _, clusterID := range item.OwnedClusters {
			_, err := service.db.Exec(`insert into owned_clusters
				(customer_id, cluster_id)
				values
				($1, $2)`,
				item.ID, clusterID)
			if err != nil {
				t.Fatal(err)
				t.Fail()
			}
		}
	}

	result, err := service.List(nil)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}

	if result.Size != expected.Size {
		t.Fatalf("Expected number of items to be %d instead got %d",
			result.Size, expected.Size)
	}

	args := &ListArguments{Page: 0, Size: 1}

	result, err = service.List(args)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}

	if result.Size != args.Size {
		t.Fatalf("Expected number of items to be %d instead got %d",
			args.Size, result.Size)
	}
	if result.Items[0].Name != "test-name0" {
		t.Fatalf("Expected result name to be %s instead got %s",
			expected.Items[0].Name, result.Items[0].Name)
	}
	if len(result.Items[0].OwnedClusters) != len(expected.Items[0].OwnedClusters) {
		t.Fatalf("Expected number of clusters to be %d instead got %d",
			len(result.Items[0].OwnedClusters), len(expected.Items[0].OwnedClusters))
	}
}

func deleteAll() {
	service.db.Exec("delete from owned_clusters")
	service.db.Exec("delete from customers")
}
