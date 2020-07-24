package main

import (
	// "log"

	"net"

	// "fmt"
	// "github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

// The non secure server is basically a similar construct as a TLS server, just with less parameters to define for TLS
// TODO to re-use the server parameters that are provided when go run is done
func startAvr(launchWithChannel bool) {
	// build port config to launch the TCP socket
	configPortTcp := ":2498"
	logger.Info().Msg("TCP Service Started non-securely on port " + configPortTcp)

	// Construct TCP server - TODO make naming convention the same over the  two files
	tcpServer, err := net.Listen("tcp", configPortTcp)
	if err != nil {
		// fmt.Println(err)
		logger.Error().Err(err).Msg("Error setting up TCP server")
		return
	}
	defer tcpServer.Close()

	// The famous loop to prevent the process from closing
	// determine to use a channel or not
	if launchWithChannel {
		// Build channel for logging service
		logChannel := make(chan string)
		if logToFile {
			go writeToDataStoreOverChannel(logChannel)
		} else if logForward && testLocal {
			go startConnection(localHost, "", "", logChannel)
		} else if logForward && testLocal == false && connectWithCredentials {
			go startConnection(remoteHost, userName, password, logChannel)
		} else {
			go startConnection(remoteHost, "", "", logChannel)
		}
		for {
			c, err := tcpServer.Accept()
			if err != nil {
				// fmt.Println(err)
				logger.Error().Err(err).Msg("Error during TCP handshaker")
				return
			}
			// When we accept a request we want a special GoRoutine to thread off the management of the buffer and data in it
			// See tcphandler.go for the way in which it handles this
			go handleTCPConnectionLogChannel(c, logChannel)
		}
	} else {
		for {
			c, err := tcpServer.Accept()
			if err != nil {
				// fmt.Println(err)
				logger.Error().Err(err).Msg("Error during TCP handshaker")
				return
			}
			// When we accept a request we want a special GoRoutine to thread off the management of the buffer and data in it
			// See tcphandler.go for the way in which it handles this
			go handleTCPConnection(c)
		}
	}
}
