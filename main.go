package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const googleMapsAPIURL = "https://maps.googleapis.com/maps/api/geocode/json"

func main() {
	// Open the input CSV file
	file, err := os.Open("input.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a reader for the CSV file
	reader := csv.NewReader(file)

	// Read the CSV headers
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}

	// Add a new header for the validation result column
	headers = append(headers, "Valid US/Canada")

	// Open the output CSV file
	outputFile, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Create a writer for the output CSV file
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write the headers to the output CSV file
	err = writer.Write(headers)
	if err != nil {
		panic(err)
	}

	// Process each row in the CSV file
	for {
		// Read the next row
		row, err := reader.Read()
		if err != nil {
			break // Reached end of file
		}

		// Get the address value from the row (assuming it's in the first column)
		address := row[0]

		// Call the Google Maps Geocoding API to validate the address
		isValid := validateAddress(address)

		// Append the validation result to the row
		row = append(row, isValid)

		// Write the row to the output CSV file
		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
}

// Calls the Google Maps Geocoding API to validate the address
func validateAddress(address string) string {
	// Prepare the request URL
	params := url.Values{}
	params.Set("address", address)
	params.Set("key", "YOUR_API_KEY") // Replace with your actual API key
	requestURL := fmt.Sprintf("%s?%s", googleMapsAPIURL, params.Encode())

	// Make the HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the response JSON
	// You'll need to use your preferred JSON decoding library here.
	// For simplicity, let's assume the response has a "results" field and a "formatted_address" field.
	// You might need to adjust this code depending on the structure of the actual response.
	// I haven't played with our Google API key yet so this is just example code. When I get an example of an 'official' response
	// I'll update it to work with that.
	
	// Here, we assume the response looks something like this:
	// {
	//   "results": [
	//     {
	//       "formatted_address": "123 Main St, City, State, Country"
	//     }
	//   ]
	// }
	// You can use a more robust JSON decoding library like "encoding/json" if needed.
	valid := false
	if strings.Contains(string(body), "United States") || strings.Contains(string(body), "Canada") {
		valid = true
	}

	// Return "Y" if the address is valid, otherwise "N"
	if valid {
		return "Y"
	} else {
		return "N"
	}
}
