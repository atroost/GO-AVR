//Required packages, net to setup a socket, java.io to serialize a message into a serialized java object.
const net = require('net');
const tls = require('tls');
const io = require('java.io');

// Define testconfiguration
const testConfig = require('../config/testConfig')
let   testExecutionSecure = 0;


// Initialize objectstreams to serialize messages
// const inputObjectStream = io.InputObjectStream;
const outputObjectStream = io.OutputObjectStream;

//create a test message and serialize it into something the AVR server understands. TODO multiple messages.

// Create the actual function, we want to run sequences of a singleTest
module.exports = {
    secureTestAvr: function (environment) {
        const environmentSelector = environment + testConfig.testTargetBaseUrl
        function sequentialMessageTesterSecure() {
            function singleMessageTestSecure() {
                // Initialize client
                const secureAvrClient = new tls.TLSSocket();

                // setup connection to configured server
                secureAvrClient.connect(testConfig.testTargetPortSecure, environmentSelector, testConfig.options, () => {
                    const today = new Date();
                    const hours = today.toLocaleTimeString();
                    const date = today.toLocaleDateString('en-GB');
                    const avrMessage = `${testConfig.testNameSecure+testExecutionSecure}|${testConfig.subevt}|${testConfig.tan}|${testConfig.ip}|${testConfig.subscriberId}|${testConfig.viewerId}|${hours + testConfig.timeZone + date }|${testConfig.interfaceVersion}|${testConfig.stbName}|${testConfig.testData}|`
                    const avrBuffer = outputObjectStream.writeObject(avrMessage);
                    console.log(today)
                    console.log(`Connected to AVR server over HTTPS at ${environmentSelector}:${testConfig.testTargetPortSecure}`)
                    
                    // Send data to connected server
                    secureAvrClient.write(avrBuffer);
                    console.log(`Wrote ${avrBuffer} securely to ${environmentSelector}`);

                    //When server signals the end of the message close the connection.
                    secureAvrClient.on('close', function () {
                        console.log(`Connection closed, attempt ${testExecutionSecure} completed`);
                        if(testExecutionSecure === testConfig.testAttempts){
                            console.log(`Secure test completed`)
                        }
                    });
                });
                // Handle errors messages
                secureAvrClient.on('error', (error) => {
                    console.log('Something went wrong in the HTTPS connection to the server');
                    console.log(error);
                });
                // Iterate number of connections.
                testExecutionSecure += 1;
                if (testExecutionSecure < testConfig.testAttempts) {
                    sequentialMessageTesterSecure();
                } 
            }
            setTimeout(singleMessageTestSecure, testConfig.timeOutValue);
        }
        // // Call the function to get started.
        sequentialMessageTesterSecure()
    }
}
