package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	broker := os.Getenv("BROKER")
	topic := os.Getenv("TOPIC")
	clientID := os.Getenv("CLIENT_ID")

	handler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("MSG: %s\n", msg.Payload())
		token := client.Publish(topic, 0, false, fmt.Sprintf("Received message: %v", msg.Payload()))
		token.Wait()
	}

	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(handler)

	opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	<-c
}
