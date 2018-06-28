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
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os"
	"testing"
)

var service *EtcdCustomersService
var localTesting = false

func TestMain(m *testing.M) {
	var err error
	if !localTesting {
		os.Exit(0)
	}
	service, err = NewEtcdCustomersService("localhost:2379")
	if err != nil {
		fmt.Printf("Could not run tests, an error occurred while trying to connect to etcd: %s\n", err)
	}
	defer deleteAll()
	defer service.Close()
	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	deleteAll()
	customer := Customer{
		Name: "fake-customer",
	}
	result, err := service.Add(customer)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if result.Name != customer.Name {
		t.Logf("Expected customer name to be %s instead got %s", customer.Name, result.Name)
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	deleteAll()
	expected := Customer{
		ID:            "some-fake-id",
		Name:          "fake-customer",
		OwnedClusters: []string{"fake-cluster-id0", "fake-cluster-id1"},
	}

	raw, err := json.Marshal(expected)
	if err != nil {
		t.Log(err)
		t.Fail()

	}
	s := string(raw)
	ctx := context.Background()
	_, err = service.cli.Put(ctx, expected.ID, s)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	result, err := service.Get(expected.ID)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if result.Name != expected.Name {
		t.Logf("Expected customer name to be %s instead got %s", result.Name, result.Name)
		t.Fail()
	}
	if len(result.OwnedClusters) != len(expected.OwnedClusters) {
		t.Logf("Expected customer to have %d clusters instead it has %d", len(result.OwnedClusters), len(expected.OwnedClusters))
		t.Fail()
	}
}

func TestList(t *testing.T) {
	deleteAll()
	items := []*Customer{
		&Customer{
			ID:            "some-fake-id0",
			Name:          "fake-customer0",
			OwnedClusters: []string{"fake-cluster-id0", "fake-cluster-id1"},
		},
		&Customer{
			ID:            "some-fake-id1",
			Name:          "fake-customer1",
			OwnedClusters: []string{"fake-cluster-id2"},
		},
	}
	expected := CustomersList{
		Total: 2,
		Size:  2,
		Page:  1,
		Items: items,
	}
	for _, customer := range items {
		raw, err := json.Marshal(customer)
		if err != nil {
			t.Log(err)
			t.Fail()

		}
		s := string(raw)
		ctx := context.Background()
		_, err = service.cli.Put(ctx, customer.ID, s)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	list, err := service.List(nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if list.Total != expected.Total {
		t.Logf("Expected result to have %d total customers instead it has %d", expected.Total, list.Total)
		t.Fail()
	}

	args := &ListArguments{
		Page: 1,
		Size: 1,
	}

	expected = CustomersList{
		Total: 1,
		Size:  1,
		Page:  0,
		Items: []*Customer{
			&Customer{
				ID:            "some-fake-id1",
				Name:          "fake-customer1",
				OwnedClusters: []string{"fake-cluster-id2"},
			},
		},
	}

	list, err = service.List(args)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	total := 0
	for _, items := range list.Items {
		if items != nil {
			total++
		}
	}
	if int64(total) != expected.Total {
		t.Logf("Expected result to have %d total customers instead it has %d", expected.Total, total)
		t.Fail()
	}
	if list.Items[0].Name != expected.Items[0].Name {
		t.Logf("Expected result to have a customer with nams %s instead it is %s", expected.Items[0].Name, list.Items[0].Name)
		t.Fail()
	}
}

func deleteAll() error {
	_, err := service.cli.Delete(context.Background(), "", clientv3.WithPrefix())
	return err
}
