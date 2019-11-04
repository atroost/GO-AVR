// Define the testScope
const avrTest = require('./tests/AVR_testClient.js')
const avrSecureTest = require('./tests/AVR_secure_testClient.js')
const testConfig = require('./config/testConfig')

// Test secure (Secure), non secure (Non-Secure), or both (All)
let testScope = 'Non-Secure' 

// Determine environment to use test ('local'), ('test'), acceptance ('acc') or production ('prod') 
let environment = 'local'

// Create testSuite function
function testSuite(testScope) {
    switch (testScope) {
        case 'All':
            async function sequentialTester() {
                // Create delay function to allow for squential execution
                function executionDelay(ms) {
                    return new Promise(resolve => setTimeout(resolve, ms));
                }
                // Start non-secure test
                await avrTest.testAvr(environment);
                
                // Await delay for execution of test 1
                await executionDelay(testConfig.testTime);

                // Start secure tests
                await avrSecureTest.secureTestAvr(environment);
            }
            sequentialTester();
            break;
        case 'Secure':
            // Kick off secure tests
            avrSecureTest.secureTestAvr(environment);
            break;
        case 'Non-Secure':
            // Start non-secure test
            avrTest.testAvr(environment);
            break;
    }
}

// Invoke testSuite
testSuite(testScope)

