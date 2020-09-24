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
	// e.g. /stable/artifacts/x86_64/metal/pxe/kernel
	if len(fragments) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fragments[0] == ""
	strm := fragments[1]
	// fragments[2] == "artifacts"
	arch := fragments[3]
	plat := fragments[4]
	frmt := fragments[5]
	art := fragments[6]

	q := u.Query()
	sha256 := (q["sha256"] != nil)
	peek := (q["peek"] != nil)
	sig := (q["sig"] != nil)

	s, err := svr.Resolve(context.Background(), strm)
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

func gohome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "https://github.com/mcsaucy/fcos-image-redirector", http.StatusFound)
}

func main() {

	svr := server{streams.New()}
	http.HandleFunc("/stable/artifacts/", svr.artifacts)
	http.HandleFunc("/testing/artifacts/", svr.artifacts)
	http.HandleFunc("/next/artifacts/", svr.artifacts)
	http.HandleFunc("/", gohome)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
