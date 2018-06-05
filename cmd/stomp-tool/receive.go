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
	"github.com/go-stomp/stomp"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive messages from a destination",
	Long:  "Receive messages from a destination.",
	Run:   runReceive,
}

func runReceive(cmd *cobra.Command, args []string) {
	// Check mandatory arguments:
	if destinationName == "" {
		glog.Errorf("The argument 'destination' is mandatory")
		return
	}

	// Connect to the messaging service:
	connection, err := connect()
	if err != nil {
		glog.Errorf(
			"Can't connect to STOMP broker at host '%s' and port %d: %s",
			brokerHost,
			brokerPort,
			err.Error(),
		)
		return
	}
	defer connection.Disconnect()
	glog.Errorf(
		"Connected to STOMP broker at host '%s' and port %d",
		brokerHost,
		brokerPort,
	)

	// Receive messages:
	subscription, err := connection.Subscribe(destinationName, stomp.AckAuto)
	if err != nil {
		glog.Errorf(
			"Can't subscribe to destination '%s': %s",
			destinationName,
			err.Error(),
		)
		return
	}
	glog.Infof(
		"Subscribed to destination '%s'",
		destinationName,
	)

	// Wait for messages:
	for message := range subscription.C {
		if message.Err != nil {
			glog.Errorf(
				"Received error message: %s",
				message.Err.Error(),
			)
			break
		}
		glog.Infof(
			"Received message from destination '%s':\n%s",
			destinationName,
			message.Body,
		)
	}
}
