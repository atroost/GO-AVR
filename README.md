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
* mqttforwarder.go - contains a forwarder module that can be used to write logrows to an MQTT endpoint;

For legacy reasons a testsuite is created in Node.js and can be found under the testClient folder. By opening the testClient in an editor like visual studio or other it should provide guidance on how to use the testclient itself. In general it aims to:
* provide a way to set the amount of testmessages to send;
* provide a way to test concurrency or sequential on a given amount of messages (and a given timimg)
* provide a way to test both secure and insecure interfaces on the receiving server (using node.js built in public cert set)

The aim for the testclient and split of the responsibility in Go files is to ensure  distributed development and easier to read go files.

### Run code
To run the code including a subfiles it's possible to provide the server arguments directly to the application (including containers if container runtime is used). A zerolog logger is used to perform overall logmanagement within the application for server logs. The standard logger of Go is used for handling writes to external files because of formatting. When launching the server loglevels can be set using the loglevel as an input parameter. When no input parameter is provided for the zerolog, the application falls back to loglevel "info".
```
go run . <port> <loglevel>
```

Using the included Dockerfile it's possible to run the container (either detached or not). To do so use the included dockerfiles per environment. The most important thing is to ensure that you have the private key and PEM file for the domain on which the server will be active

### Build docker image for production
Create folder for the production certificates and call this folder certs-prod. Place the pem and the keyfile in this directory and ensure that the naming convention is server.pem or server.key. The dockerfile will ensure that only the applicable environment certificates are copied. Check needed is if you want to expose port 2498 or port 2499.
```
 docker build -t <avr_image>:version .
```
### Build docker image for SIT
Create folder for the SIT certificates and call this folder certs-sit. Place the pem and the keyfile in this directory and ensure that the naming convention is server.pem or server.key. The dockerfile will ensure that only the applicable environment certificates are copied. Check needed is if you want to expose port 2498 or port 2499.
```
 docker build -f Dockerfile_SIT -t <avr_image>:version .
```
Run the associated created image (this case OSX forces to explicitely set the port for docker) As an input parameter we also put the port into  the container runtime (2498 for non-secure / 2499 for secure)

### Run docker image
```
 docker run -p 2498:2498 <avr_image_reference> 2498/2499 tail -f /dev/null
```

This will launch the command line tailing interface of the container itself. Use the testclient to see logfiles being created in the container

### Run test harnass
```
 npm start
```
### Exporting the docker to Nexus
When you've succesfully tested the container to run in standalone mode, it's time to upload the image to a nexus that will be used for deployment. In the scenario below we'll use a Kubernetes platform that has Nexus integrated into it. For other uploads, check your provider's nexus location. 

ensure that you are logged in to the nexus platform
```
 docker login -u <username> -p password <endpoint for nexus>
```
create a nexus tag for the build you want to push to the endpoint, in this example we'll use avr_go as the name of the deployment
```
 docker tag avr_go <endpoint>/avr_go:<version>
```
Push the container to the nexus endpoint, in this example we'll use avr_go as the name of the deployment and we'll re-use the tag provided in the previous step.
```
 docker push <endpoint>/avr_go:<version>
```
From here you can use regular kubernetes deployment files to launch your application (see K8s for examples).