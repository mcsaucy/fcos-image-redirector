package main

import (
	"context"
	"fcos-image-redirector/streams"
	"fmt"
	"log"
)

func main() {
	s, err := streams.New().Resolve(context.Background(), "stable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("parsed streams: %#v\n", s)
}
