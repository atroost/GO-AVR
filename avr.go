package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"os"
	"strings"
	"fmt"
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
	// Create unique ID to be able to generate unique logfile names
	//The string representation of a UUID consists of 32 hexadecimal digits displayed in 5 groups but NOT separated by hyphens.
	b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    uuid := fmt.Sprintf("%x%x%x%x%x",
        b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    fmt.Printf("Unique id of AVR server: %s\n", uuid)
	
	// create unique log
	f, err := os.OpenFile(uuid + ".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	if !useSecureServer {
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