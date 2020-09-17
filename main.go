package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func streamjson(stream string) ([]byte, error) {
	url := "https://builds.coreos.fedoraproject.org/streams/" + stream + ".json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %v: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-OK status when fetching %v: %v", url, resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of %v: %w", url, err)
	}
	return bodyBytes, nil
}

func main() {
	j, err := streamjson("stable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
}
