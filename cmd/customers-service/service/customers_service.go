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

package service

// CustomersService is an interface exposing a set of operations required for
// running and operating the customers of the Openshift Dedicated Portal.
type CustomersService interface {

	// List returns a pointer to CustomerList or error in case some error occurred.
	// If list arguments are provided list will return the intended customers list.
	// If nil is supplied list will return all customers.
	List(args *ListArguments) (*CustomersList, error)

	// Add creates a customer and returns the newly created customer or error
	// in case some error occurred.
	// It receives a Customer object with its Name and (possibly) OwnedClusters,
	// and creates a new Customer based on the supplied Customer parameter.
	Add(customer Customer) (*Customer, error)

	// Get returns a pointer to customer with id supplied or error if an
	// error occurred.
	// If no such customer exist Get returns nil pointer and nil error.
	Get(id string) (*Customer, error)

	// Close closes the service.
	Close() error
}

// ListArguments are arguments relevant for listing objects
type ListArguments struct {
	Page int64
	Size int64
}
