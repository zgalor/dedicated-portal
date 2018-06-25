package main

import (
	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

type Notifier struct {
	stopCh <-chan struct{}
}

func NewNotifier(stopCh <-chan struct{}) *Notifier {
	notifier := new(Notifier)
	notifier.stopCh = stopCh
	return notifier
}

func (n Notifier) SendNotification(message, destinationName string) error {
	var c client.Connection
	var err error

	// Set the clients variables before we can open it.
	c, err = stomp.NewConnection(&stomp.ConnectionBuilder{
		// Global options:
		BrokerHost:   "messaging-service.dedicated-portal.svc",
		BrokerPort:   61613,
		UserName:     "clusters-service",
		UserPassword: "redhat123",
		UseTLS:       true,
		InsecureTLS:  true,
	})
	if err != nil {
		return err
	}
	// Connect to the messaging service:
	err = c.Open()
	if err != nil {
		return err
	}
	defer c.Close()
	m := client.Message{
		ContentType: "text/plain",
		Body:        message,
	}
	err = c.Publish(m, destinationName)
	if err != nil {
		return err
	}
	return nil

}
