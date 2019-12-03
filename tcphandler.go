package main

import (
	"bufio"
	// "fmt"
	"io"
	"net"
	logger "github.com/rs/zerolog/log"
	"strings"
	// "regexp"
	// "unicode"
)

// Create alternative handler for serving TCP connections
func handleTCPConnectionLogChannel(c net.Conn, logChannel chan string) {
	logger.Debug().Msg("Handling incoming socket: "+ c.RemoteAddr().String())

	// Create buffer to which data can be written
	// Ensure buffer can be read with bufio package from Golang
	dataBuffer := make([]byte, 4096)
	bufferReader := bufio.NewReader(c)
	
	// Create loop to get the size of the message and thereby outputting the entire message
	for {
		// Start preparing capability to read from the generated buffer
		bufferByteReader, err := bufferReader.ReadByte()
		if err != nil {
			// fmt.Println("Error hit in bytereader: ",err)
			logger.Error().Err(err).Msg("Error hit in bytereader function")
			return
		}

		// Based on the amount of data buffered in the dataBuffer buffer we can get data from it
		// For some reason this function only works with bufferByteReader active
		dataInBuffer := bufferReader.Buffered()
		// fmt.Println("Buffersize is ", dataInBuffer)

		// read the full message, or return an error
		readBytes, err := io.ReadFull(bufferReader, dataBuffer[:int(dataInBuffer)])
		if err != nil {
			// fmt.Println("Error hit in readbyte function: ",err)
			logger.Error().Err(err).Msg("Error hit in readbyte function")
			return
		}

		// Convert buffer to string if there is more than 0 bytes available to convert
		// fmt.Println(dataBuffer[:int(dataInBuffer)])
		if bufferByteReader > 0 && readBytes > 0 {
			// Convert data to string for logging purposes
			bytesConvertedToString := string(dataBuffer[:int(dataInBuffer)])

			// Validate message index based on available log messages
			messageIndex := strings.IndexAny(bytesConvertedToString, "ABCDEGHMNOPQRSTUV")
			logger.Trace().Str("Index", string(messageIndex)).Msg("Calculated messageindex for message")

			// Use index to get full message
			indexedMessage := string(dataBuffer[messageIndex:int(dataInBuffer)])
			
			// OR Create regexp to ensure all unwanted characters are stripped
			// avrRegex, err := regexp.Compile("[^a-zA-Z0-9|\n ]+")
			// if err != nil {
			// 	logger.Trace().Err(err).Msg("regular expression error")
			// }
			// processedAvrString := avrRegex.ReplaceAllString(bytesConvertedToString, "")
			
			// Strip of all unwanted characters from the regular expression
			// stripProcessedString := strings.TrimLeft(processedAvrString, "abcdefghikjlmnopqrstuvwxyz1234567890{}?!@#$%^&*()[]123456789�")			
			
			// Optional loglines for debugging
			logger.Trace().Str("bytesize", string(bufferByteReader)).Msg("BufferBytereader calculated")

			logChannel <- indexedMessage
			c.Close()
			return
		}
		return
	}
	
}

// Create alternative handler for serving TCP connections
func handleTCPConnection(c net.Conn) {
	logger.Debug().Msg("Handling incoming socket: "+ c.RemoteAddr().String())

	// Create buffer to which data can be written
	// Ensure buffer can be read with bufio package from Golang
	dataBuffer := make([]byte, 4096)
	bufferReader := bufio.NewReader(c)
	
	// Create loop to get the size of the message and thereby outputting the entire message
	for {
		// Start preparing capability to read from the generated buffer
		bufferByteReader, err := bufferReader.ReadByte()
		if err != nil {
			// fmt.Println("Error hit in bytereader: ",err)
			logger.Error().Err(err).Msg("Error hit in bytereader function")
			return
		}

		// Based on the amount of data buffered in the dataBuffer buffer we can get data from it
		// For some reason this function only works with bufferByteReader active
		dataInBuffer := bufferReader.Buffered()
		// fmt.Println("Buffersize is ", dataInBuffer)

		// read the full message, or return an error
		readBytes, err := io.ReadFull(bufferReader, dataBuffer[:int(dataInBuffer)])
		if err != nil {
			// fmt.Println("Error hit in readbyte function: ",err)
			logger.Error().Err(err).Msg("Error hit in readbyte function")
			return
		}

		// Convert buffer to string if there is more than 0 bytes available to convert
		if bufferByteReader > 0 && readBytes > 0 {
			// Convert data to string for logging purposes
			bytesConvertedToString := string(dataBuffer[:int(dataInBuffer)])
			// bytesConvertedToUtf := strings.ToValidUTF8(bytesConvertedToString, "")
			stringStrippedOfUnicode := strings.TrimLeft(bytesConvertedToString, "\ufffd\u0000\u0005t\u0000\\ufff\u0003\ufffd\u0001A\u00001\u00002&abcdefghikjlmnopqrstuvw1234567890{}!@#$%^&*()")

			// Optional loglines for debugging
			// fmt.Println("BufferBytereader size: ", bufferByteReader)
			logger.Trace().Msg("BufferBytereader calculated: " + string(bufferByteReader))

			// Create goroutine that stores data to different functions
			go writeToDataStore(stringStrippedOfUnicode)
			// go writeToRemoteEndpoint(stringConversion)
			c.Close()
			return
		}
		c.Close()
		return
	}
}