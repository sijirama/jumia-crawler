package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var key = "fb9d8506-af67-4416-9f3a-a44f1a236d81"

type FakeHeadersResponse struct {
	Result []map[string]string `json:"result"`
}

func RandomHeaders(headersList []map[string]string) map[string]string {
	randomIndex := rand.Intn(len(headersList))
	return headersList[randomIndex]
}

func GetHeadersList() []map[string]string {
	scrapeopsAPIEndpoint := "https://headers.scrapeops.io/v1/browser-headers?api_key=" + key + "&num_results=20"
	req, _ := http.NewRequest("GET", scrapeopsAPIEndpoint, nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var fakeHeadersResponse FakeHeadersResponse
		json.NewDecoder(resp.Body).Decode(&fakeHeadersResponse)
		return fakeHeadersResponse.Result
	}
	var emptySlice []map[string]string
	fmt.Println(emptySlice)
	return emptySlice
}

func GetRandomHeaders() map[string]string {
	headersList := GetHeadersList()             // Get the list of headers
	randomHeaders := RandomHeaders(headersList) // Get the random headers

	// Create a new map to hold the combined headers
	finalHeaders := make(map[string]string)

	// Set the specified headers
	finalHeaders["Sec-Fetch-Site"] = "same-origin"
	finalHeaders["Sec-Fetch-Mode"] = "navigate"
	finalHeaders["Sec-Fetch-User"] = "?1"
	finalHeaders["Upgrade-Insecure-Requests"] = "1"

	// Append random headers to the final map
	for key, value := range randomHeaders {
		finalHeaders[key] = value
	}

	return finalHeaders // Return the combined headers
}
