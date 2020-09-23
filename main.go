package main

import (
	"context"
	"fmt"
	"github.com/mcsaucy/fcos-image-redirector/streams"
	"log"
	"net/http"
	"strings"
)

type server struct {
	streams.Resolver
}

func (svr *server) artifacts(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	fragments := strings.Split(u.Path, "/")
	// e.g. /artifacts/x86_64/metal/pxe/kernel
	if len(fragments) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fragments[0] == ""
	// fragments[1] == "artifacts"
	arch := fragments[2]
	plat := fragments[3]
	frmt := fragments[4]
	art := fragments[5]

	q := u.Query()
	sha256 := (q["sha256"] != nil)
	peek := (q["peek"] != nil)
	sig := (q["sig"] != nil)

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

	if sha256 {
		fmt.Fprintf(w, "%v\n", res.Sha256)
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
