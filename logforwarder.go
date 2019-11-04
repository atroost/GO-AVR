package main 

import ( 
	"net/http"
	"bytes"
)

// Create function to forward data to a remote location
func writeToRemoteEndpoint(dataToWrite string) {
	// Create byte representation of data to forward
	jsonStr := []byte(dataToWrite)

	// Currently Hookbin is used as an endpoint for sending dummy data to, we can add other locations here later
	// Set request structure towards hookbin && set headers
	req, err := http.NewRequest("POST", "https://hookb.in/OebqpkVda0tMpwV2d1X0", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	
	// Create client that can post to hookbin
	client := &http.Client{}
	
	// Function to create a request towards hookbin + error handling
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	
	defer resp.Body.Close()
	return
}