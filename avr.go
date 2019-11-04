package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"os"
	"strings"
)

func handleConnection(c net.Conn) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
	}
	c.Close()
}

// only needed below for sample processing

func main() {
	// determine whether to use a secure server or not
	useSecureServer := false
	// create log
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	if useSecureServer {
		startAvrNoCert()
	} else {
		// check certs
		cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		// verify port is part of command
		arguments := os.Args
		if len(arguments) == 1 {
			log.Printf("Please provide a port number!")
			return
		}

		log.Printf("Launching server...")
		// define port
		PORT := ":" + arguments[1]
		// listen on all interfaces
		config := tls.Config{Certificates: []tls.Certificate{cert}}
		config.Rand = rand.Reader
		ln, err := tls.Listen("tcp", PORT, &config)
		if err != nil {
			log.Println(err)
			return
		}

		// run loop forever (or until ctrl-c)
		for {
			c, err := ln.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			go handleConnection(c)
		}
	}
}