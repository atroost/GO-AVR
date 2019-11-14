//Required packages, net to setup a socket, java.io to serialize a message into a serialized java object.
const net = require('net');
const io = require('java.io');

// Define testconfiguration
const testConfig = require('../config/testConfig')
let   testExecution = 0;
// let avrTestComplete = false

// Initialize objectstreams to serialize messages
// const inputObjectStream = io.InputObjectStream;
const outputObjectStream = io.OutputObjectStream;

//create a test message and serialize it into something the AVR server understands. TODO multiple messages, pull messages from central config.

// Create the actual function, we want to run sequences of a singleTest
module.exports = {
    testAvr: function (environment) {
        let environmentSelector
        if (environment === 'local') {
            environmentSelector = 'localhost';
        } else {
            environmentSelector = environment + testConfig.testTargetBaseUrl 
        }
        function sequentialMessageTester() {
            function singleMessageTest() {
                // setup connection to configured server
                const avrClient = new net.Socket();
                avrClient.connect(testConfig.testTargetPort, environmentSelector, () => {
                    const today = new Date();
                    const hours = today.toLocaleTimeString();
                    const date = today.toLocaleDateString('en-GB');
                    const avrMessage = `${testConfig.testName+testExecution}|${testConfig.subevt}|${testConfig.tan}|${testConfig.ip}|${testConfig.subscriberId}|${testConfig.viewerId}|${hours + testConfig.timeZone + date }|${testConfig.interfaceVersion}|${testConfig.stbName}|${testConfig.testData}|`
                    // const normalizer = outputObjectStream.normalize(avrMessage, string)
                    const avrBuffer = outputObjectStream.writeObject(avrMessage);
                    console.log(today)
                    console.log(avrBuffer)
                    console.log("AVR message: " + avrMessage)
                    console.log(`Connected to AVR server over HTTP at ${environmentSelector}:${testConfig.testTargetPort}`);
                    
                    // To test if receiver and sender are seeing the same buffers check size.
                    function byteCount(s) {
                        return encodeURI(s).split(/%..|./).length - 1;
                    }
                    // byteCount(avrMessage)
                    console.log(`Length of AVR message is ${byteCount(avrMessage)}`);

                    // Send data to connected server
                    avrClient.write(avrMessage);
                    console.log(`Wrote ${avrMessage} non-securely to ${environmentSelector}`);

                    //When server signals the end of the message close the connection.
                    avrClient.on('close', function () {
                        console.log(`Connection closed, attempt ${testExecution} completed`);
                        if(testExecution === testConfig.testAttempts){
                            console.log(`Non secure test completed`)
                            // testConfig.testComplete = true 
                        }
                    });
                });
                // Handle errors messages
                avrClient.on('error', (error) => {
                    console.log('Something went wrong in the http connection to the server');
                    console.log(error);
                });
                // Iterate number of tests.
                testExecution += 1;
                if (testExecution < testConfig.testAttempts) {
                    sequentialMessageTester();
                } 
            }
            setTimeout(singleMessageTest, testConfig.timeOutValue);
        }
        // Call the function to get started.
        sequentialMessageTester();
    }
};

