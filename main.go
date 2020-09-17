package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://builds.coreos.fedoraproject.org/streams/stable.json")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to fetch streams json: %w", resp)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bodyBytes))
}
