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
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/segmentio/ksuid"
)

// EtcdCustomersService is a struct implementing the customer service interface,
// backed by an etcd cluster.
type EtcdCustomersService struct {
	cli *clientv3.Client
}

// NewEtcdCustomersService is a constructor for the EtcdCustomersService struct.
func NewEtcdCustomersService(etcdEndpoint string) (service *EtcdCustomersService, err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	service = new(EtcdCustomersService)
	service.cli = cli
	return service, nil
}

// Close closes the etcd customers service client.
func (service *EtcdCustomersService) Close() {
	service.cli.Close()
}

// Add adds a single customer to etcd cluster.
func (service *EtcdCustomersService) Add(customer Customer) (*Customer, error) {
	// generate customer id.
	id, err := ksuid.NewRandom()

	if err != nil {
		return nil, err
	}

	// result is the new customer inserted into etcd.
	result := Customer{
		ID:   id.String(),
		Name: customer.Name,
	}

	if customer.OwnedClusters == nil {
		result.OwnedClusters = make([]string, 0)
	} else {
		result.OwnedClusters = customer.OwnedClusters
	}

	// the resulting Customer is then marshal to []byte and converted to string -
	// since etcd key-value pairs are strings ONLY.
	// The key used as reference for Customer objects is the customer id.
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	s := string(raw)
	ctx := context.Background()
	_, err = service.cli.Put(ctx, result.ID, s)

	if err != nil {
		return nil, err
	}
	return &result, err
}

// Get retrieves a single customer from etcd cluster
func (service *EtcdCustomersService) Get(id string) (*Customer, error) {
	// retrieve customer object by it's id.
	response, err := service.cli.Get(context.Background(), id)
	if err != nil {
		// If could not find customer matching such id return nil pointer and nil error.
		return nil, nil
	}

	// We expect only one Customer per ID since ID's are unique.
	if response.Count != 1 {
		return nil, fmt.Errorf("key %s should contain a single value, instead it contains %d values", id, response.Count)
	}

	result := new(Customer)
	for _, ev := range response.Kvs {
		json.Unmarshal(ev.Value, result)
	}
	return result, nil
}

// List retrieves a list of current customers stored in datastore.
func (service *EtcdCustomersService) List(args *ListArguments) (*CustomersList, error) {
	// We get all Customer objects by querying etcd for object with empty-prefix.
	response, err := service.cli.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	kvs := response.Kvs
	customerList, err := service.paginateCustomers(args, response.Count, kvs)

	if err != nil {
		return nil, err
	}
	return customerList, nil
}

// paginateCustomers returns a *CustomersList representing a single page of
// Customers information.
func (service *EtcdCustomersService) paginateCustomers(args *ListArguments, total int64, keyValues []*mvccpb.KeyValue) (*CustomersList, error) {
	// if no list arguments specified - get all customers.
	var page int64
	var size int64
	var firstIndex int64
	var lastIndex int64

	if args == nil {
		page = 1
		size = total
		firstIndex = 0
		lastIndex = total
	} else {
		page = args.Page
		size = args.Size
		firstIndex = size * page
		lastIndex = size * (page + 1)
	}

	items, err := service.collectCustomerItems(size, page, total, firstIndex, lastIndex, keyValues)
	if err != nil {
		return nil, err
	}

	// Return the customer list for requested page.
	result := CustomersList{
		Items: items,
		Page:  page,
		Size:  int64(len(items)),
		Total: total,
	}
	return &result, nil
}

func (service *EtcdCustomersService) collectCustomerItems(size, page, total, firstIndex, lastIndex int64, keyValues []*mvccpb.KeyValue) ([]*Customer, error) {
	items := make([]*Customer, size)
	for i := firstIndex; i < lastIndex && i < total; i++ {
		err := json.Unmarshal(keyValues[i].Value, &items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (service *EtcdCustomersService) deleteAll() error {
	_, err := service.cli.Delete(context.Background(), "", clientv3.WithPrefix())
	return err
}
