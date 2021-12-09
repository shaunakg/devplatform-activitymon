// This service is designed to run on a server as a monitoring script.
// It will check the http://localhost/healthz endpoint, which will return a JSON response.

// If the JSON response indicates activity in the last 30 minutes, the script will exit.
// Otherwise (or if the request fails), the script will make a HTTP GET request to https://srg-devplatform-backend-k7em4uqi4a-km.a.run.app/self/stop

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// healthzResponse is a struct that will be unmarshalled from the JSON response from the healthz endpoint.
	// It contains a status and lastHeartbeat key.
	// The status key is a string, and lastHeartbeat is a unix timestamp.
	type healthzResponse struct {
		Status        string `json:"status"`
		LastHeartbeat int64  `json:"lastHeartbeat"`
	}

	// Make a HTTP get request to the healthz endpoint.
	resp, err := http.Get("http://localhost/healthz")
	if err != nil {

		// If the request fails, make a HTTP get request to the stop endpoint.
		resp, err = http.Get("https://srg-devplatform-backend-k7em4uqi4a-km.a.run.app/self/stop")
		if err != nil {
			log.Fatal(resp)
			log.Fatal("Error making stop request after healthz failure: ", err)
		}

		os.Exit(0)

	}

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the response body into the healthzResponse struct.
	var healthz healthzResponse
	err = json.Unmarshal(body, &healthz)
	if err != nil {
		log.Fatal(err)
	}

	// If the timestamp is less than 30 minutes ago, exit.
	if healthz.LastHeartbeat > (time.Now().Unix() - 1800) {
		fmt.Println("Activity in the last 30 minutes, timestamp: ", healthz.LastHeartbeat)
		os.Exit(0)
	}

	// If the timestamp is more than 30 minutes ago, make a HTTP get request to the self/stop endpoint.
	resp, err = http.Get("https://srg-devplatform-backend-k7em4uqi4a-km.a.run.app/self/stop")
	if err != nil {
		log.Fatal(err)
	}

	// If the response is not 200, fatally log
	if resp.StatusCode != 200 {
		// Log the response body.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(body))
	}

	// If the response is 200, exit.
	fmt.Println("This instance should be shutting down now.")
	os.Exit(0)

}
