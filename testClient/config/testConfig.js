// Configuration data for running various tests

const testConfig = {}

// Server configuration
testConfig.testTargetBaseUrl = '.log.stb.itvonline.nl';
testConfig.testName = 'avrTester';
testConfig.testNameSecure = 'secureAvrTester';
testConfig.testTargetPort = 2498;
testConfig.testTargetPortSecure = 2499;
testConfig.options = {
  rejectUnauthorized: true,
  enableTrace: true,
};
testConfig.testComplete = false;
testConfig.testSecureComplete = false;

// Test configuration
testConfig.testAttempts = 1;
testConfig.timeOutValue = 100;
testConfig.testTime = testConfig.testAttempts*testConfig.timeOutValue;

// Message configuration
testConfig.subevt = 'SUBEVT'
testConfig.tan = '00012345678900'
testConfig.ip = '192.168.2.2'
testConfig.subscriberId = '2'
testConfig.viewerId = '1'
testConfig.timeZone = ' GMT ' 
testConfig.interfaceVersion = 'ver=HE31.V01'
testConfig.stbName = 'Living Room'
testConfig.testData = '15616641|CANVAS|Access Hollywood|02/02/2004|09:00:00|10:59:59|STOPPED|AGR7342982|HD|'    
// const avrMessage = `${testName+testExecution}|${testConfig.subevt}|${testConfig.tan}|${testConfig.ip}|${testConfig.subscriberId}|${testConfig.viewerId}|${hours + testConfig.timeZone + date }|${testConfig.interfaceVersion}|${testConfig.stbName}|${testConfig.testData}|`
 

module.exports = testConfig