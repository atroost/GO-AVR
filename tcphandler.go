package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Create alternative handler for serving TCP connections
func handleTCPConnection(c net.Conn) {
	fmt.Printf("Handling socket from %s\n", c.RemoteAddr().String())

	// Create buffer to which data can be written
	// Ensure buffer can be read with bufio package from Golang
	dataBuffer := make([]byte, 4096)
	bufferReader := bufio.NewReader(c)

	// Create loop to get the size of the message and thereby outputting the entire message
	for {
		// Start preparing capability to read from the generated buffer
		bufferByteReader, err := bufferReader.ReadByte()
		if err != nil {
			fmt.Println(err)
		}

		// Based on the amount of data buffered in the dataBuffer buffer we can get data from it
		// For some reason this function only works with bufferByteReader active
		dataInBuffer := bufferReader.Buffered()
		fmt.Println("Buffersize is ", dataInBuffer)

		// read the full message, or return an error
		readBytes, err := io.ReadFull(bufferReader, dataBuffer[:int(dataInBuffer)])
		if err != nil {
			fmt.Println(err)
		}

		// Convert buffer to string if there is more than 0 bytes available to convert
		if bufferByteReader > 0 && readBytes > 0 {
			// Convert data to string for logging purposes
			stringConversion := string(dataBuffer[:int(dataInBuffer)])

			// Optional loglines for debugging
			fmt.Println("BufferBytereader size: ", bufferByteReader)

			// Create goroutine that stores data to different functions
			go writeToDataStore(stringConversion)
			go writeToRemoteEndpoint(stringConversion)
			return
		}
	}
}