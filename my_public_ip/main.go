package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://ifconfig.me"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch IP address: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("Public IP Address: %s\n", string(body))
}
