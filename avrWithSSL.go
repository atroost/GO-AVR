package main

import (
	// "fmt"
	"crypto/rand"
	"crypto/tls"
	// "log"
	logger "github.com/rs/zerolog/log"
)

// The non secure server is basically a similar construct as a TLS server, just with less parameters to define for TLS
// TODO to re-use the server parameters that are provided when go run is done
func startAvrSecure(launchWithChannel bool) {
	// build port config to launch the TCP socket
	configPortSecureTcp := ":2499"
	// fmt.Println("** TCP Service Started securely on Port:", configPortSecureTcp, " **")
	logger.Info().Msg("TCP Service Started securely on port " + configPortSecureTcp)
	
	// load certificates
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		// log.Fatalf("server: loadkeys: %s", err)
		logger.Fatal().Err(err).Msg("Error while loading certificates")
	}

	// Construct secure TCP server and make it compliant with KSP.
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}
	config.Rand = rand.Reader
	secureTcpServer, err := tls.Listen("tcp", configPortSecureTcp, &config)
	if err != nil {
		// fmt.Println(err)
		logger.Error().Err(err).Msg("Error setting up Secure TCP server")
		return
	}

	// The famous loop to prevent the process from closing
	// determine to use a channel or not
	if launchWithChannel {
		// Build channel for logging service
		logChannel:= make(chan string)
		go writeToDataStoreOverChannel(logChannel)
		for {
			c, err := secureTcpServer.Accept()
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
				c, err := secureTcpServer.Accept()
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

