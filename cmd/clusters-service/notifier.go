package main

import (
	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

// Notifier notify about things.
type Notifier struct {
	stopCh <-chan struct{}
}

// NewNotifier create a new notifier.
func NewNotifier(stopCh <-chan struct{}) *Notifier {
	notifier := new(Notifier)
	notifier.stopCh = stopCh
	return notifier
}

// SendNotification send a notification.
func (n Notifier) SendNotification(message, destinationName string) error {
	var c client.Connection
	var err error

	// Set the clients variables before we can open it.
	c, err = stomp.NewConnection(&client.ConnectionSpec{
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
	defer c.Close()
	m := client.Message{
		ContentType: "text/plain",
		Data: client.MessageData{
			"text": message,
		},
	}
	err = c.Publish(m, destinationName)
	if err != nil {
		return err
	}
	return nil

}
