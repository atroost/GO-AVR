package main

import (
	"fmt"
	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/rs/zerolog/log"
)

//define a function for the default message handler -
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	logger.Debug().Msg("Topic of message is: " + msg.Topic())
	logger.Debug().Msg("Payload of message is: " + string(msg.Payload()))
}

//define a function for the default message handler - connection setup
var connectionHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	logger.Info().Msg("Connection with MQTT broker (re-) established")
}

//define a function for the default message handler - connection lost
var connectivityLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, reason error) {
	logger.Error().Err(reason).Msg("Connection lost")
}

// Create function to start a persistent MQTT connection with the broker.
func startMqttConnection(brokerurl string, username string, password string, logChannel chan string) {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + brokerurl)
	opts.SetClientID(mqttClientName)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetConnectionLostHandler(connectivityLostHandler)
	opts.SetOnConnectHandler(connectionHandler)
	// opts.SetReconnectingHandler(connectionReconnectHandler)
	opts.SetAutoReconnect(true)
	logger.Info().Msg("Starting mqtt connection at " + brokerurl)

	//create and start a client using the above ClientOptions
	mqttPublisher := MQTT.NewClient(opts)
	if token := mqttPublisher.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	//Setup subscription to a channel
	if token := mqttPublisher.Subscribe(mqttSubscribeTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	} else {
		logger.Info().Msg("Subscribed to topic: " + mqttSubscribeTopic)
	}
	//Create for loop that listens continuously to data coming in over the message channel
	for {
		// Retrieve data from the channel on which data is published
		payload := <-logChannel
		logger.Debug().Msg("Receiving " + payload)
		// Publish data to mqtt
		token := mqttPublisher.Publish(mqttPublishTopic, 0, false, payload)
		token.Wait()
	}
}
