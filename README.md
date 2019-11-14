# GO-AVR
## AVR Server in GoLang

The AVR server can be launched in a secure and non-secure way. At the moment, the non-secure way can be launched at any time (even  if no launch parameters have been provided) while it's using a different method to handle the TCP connection.

The application structure is currently as follows:
* avr.go  - contains the most important main executionable that directs to the right component;
* avrWithSSL.go - contains a hardcoded secure TCP server that listens to TCP connections over the hardcoded port;
* avrWithoutSSL.go - contains a hardcoded non secure TCP server that listens to TCP connections over the hardcoded port;
* tcphandler.go - contains a TCP handler that reads data from a TCP Bytestream and converts it into strings;
* logToFile.go - contains a writer module that can be used to write logrows to whichever preferred logging utility or stdout;
* logForwarder.go - contains a forwarder module that can be used to write logrows to whichever preferred endpoint;

For legacy reasons a testsuite is created in Node.js and can be found under the testClient folder. By opening the testClient in an editor like visual studio or other it should provide guidance on how to use the testclient itself. In general it aims to:
* provide a way to set the amount of testmessages to send;
* provide a way to test concurrency or sequential on a given amount of messages (and a given timimg)
* provide a way to test both secure and insecure interfaces on the receiving server (using node.js built in public cert set)

The aim for the testclient and split of the responsibility in Go files is to ensure  distributed development and easier to read go files.

To run the code including a subfiles:
```
go run .
```

Using the included Dockerfile it's possible to run the container (either detached or not). To do so follow the upcoming steps

###Build docker image
```
 docker build -t <avr_image> .
```

Run the associated created image (this case OSX forces to explicitely set the port for docker) As an input parameter we also put the port into  the container runtime (2498 for non-secure / 2499 for secure)

###Run docker image
```
 docker run -p 2498:2498 <avr_image_reference> 2498 tail -f /dev/null
```

This will launch the command line tailing interface of the container itself. Use the testclient to see logfiles being created in the container

###Run test harnass
```
 npm start
```