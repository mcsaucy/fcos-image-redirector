package main

import (
	"context"
	"fcos-image-redirector/streams"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var (
	artifactsParser = regexp.MustCompile(`x86_64/artifacts/([[:word:]]+)/([[:word:]]+)/([[:word:]]+)$`)
)

func x86_64Artifacts(w http.ResponseWriter, r *http.Request) {
	// TODO(mcsaucy): cache this between runs.
	s, err := streams.New().Resolve(context.Background(), "stable")
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "failed to resolve stream: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO(mcsaucy): find a sexier way to do this?
	matches := artifactsParser.FindStringSubmatch(r.URL.Path)
	if len(matches) != 4 { // one per fragment + the whole string match
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	platform := matches[1]
	format := matches[2]
	artifact := matches[3]

	res := s.Architectures["x86_64"].Artifacts[platform].Formats[format][artifact]
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO(mcsaucy): have this redirect once we have a param to prevent that.
	// This present behavior is useful for development and while we need this
	// to redirect, I don't want my cURL invocation pulling down megabytes of
	// garbage each test.
	fmt.Fprintf(w, "%v\n", res.Location)
}

func main() {
	http.HandleFunc("/x86_64/artifacts/", x86_64Artifacts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
