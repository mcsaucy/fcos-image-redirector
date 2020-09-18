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
	artifactsParser = regexp.MustCompile(`artifacts/([[:word:]]+)/([[:word:]]+)/([[:word:]]+)/([[:word:]]+)$`)
)

type server struct {
	streams.Resolver
}

func (svr *server) artifacts(w http.ResponseWriter, r *http.Request) {
	// TODO(mcsaucy): find a sexier way to do this?
	matches := artifactsParser.FindStringSubmatch(r.URL.Path)
	if len(matches) != 5 { // one per fragment + the whole string match
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	arch := matches[1]
	plat := matches[2]
	frmt := matches[3]
	art := matches[4]

	peek := (r.URL.Query()["peek"] != nil)
	sig := (r.URL.Query()["sig"] != nil)

	s, err := svr.Resolve(context.Background(), "stable")
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "failed to resolve stream: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := s.Architectures[arch].Artifacts[plat].Formats[frmt][art]
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tgt := res.Location
	if sig {
		tgt = res.Signature
	}

	if peek {
		fmt.Fprintf(w, "%v\n", tgt)
		return
	}
	http.Redirect(w, r, tgt, http.StatusFound)
}

func main() {
	svr := server{streams.New()}
	http.HandleFunc("/artifacts/", svr.artifacts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
