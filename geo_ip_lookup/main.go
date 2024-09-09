package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GeoIP struct {
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

func showHelp() {
	fmt.Println("Usage: go run main.go <IP_ADDRESS>")
	fmt.Println("Example: go run main.go 8.8.8.8")
	fmt.Println()

	fmt.Println("Usage: geoip <IP_ADDRESS>")
	fmt.Println("Example: geoip 8.8.8.8")
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		showHelp()
		return
	}

	ip := os.Args[1]
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to make request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var geoIP GeoIP
	if err := json.NewDecoder(resp.Body).Decode(&geoIP); err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return
	}

	fmt.Printf("IP Address: %s\n", geoIP.Query)
	fmt.Printf("Country: %s\n", geoIP.Country)
	fmt.Printf("Region: %s\n", geoIP.RegionName)
	fmt.Printf("City: %s\n", geoIP.City)
	fmt.Printf("Latitude: %.4f, Longitude: %.4f\n", geoIP.Lat, geoIP.Lon)
	fmt.Printf("ISP: %s\n", geoIP.ISP)
	fmt.Printf("Timezone: %s\n", geoIP.Timezone)
}
