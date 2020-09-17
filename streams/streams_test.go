package streams

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestStreamsResolution(t *testing.T) {
	testdata, err := ioutil.ReadFile("testdata/test.json")
	if err != nil {
		t.Fatal("failed to load testdata")
	}

	want_url := "https://builds.coreos.fedoraproject.org/streams/test.json"
	hitTestEndpoint := false

	cli := NewTestClient(func(req *http.Request) *http.Response {
		if want_url != req.URL.String() {
			t.Errorf("Resolve(stable) requested %v; want %v", req.URL.String(), want_url)
		}
		hitTestEndpoint = true
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(testdata)),
			Header:     make(http.Header),
		}
	})

	sut := Resolver{cli: cli}
	strm, err := sut.Resolve(context.Background(), "test")
	if err != nil {
		t.Fatalf("Resolve(stable) err = %v; want nil", err)
	}

	if !hitTestEndpoint {
		t.Fatal("Resolve(stable) didn't hit the test endpoint")
	}

	type fieldFetcher func(Stream) string

	cases := map[string]struct {
		name string
		run  fieldFetcher
		want string
	}{
		"stream name": {
			run:  func(s Stream) string { return s.Name },
			want: "test",
		},
		"pxe kernel location": {
			run: func(s Stream) string {
				fmts := s.Architectures["x86_64"].Artifacts["metal"].Formats
				return fmts["pxe"]["kernel"].Location
			},
			want: "https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20200824.3.0/x86_64/fedora-coreos-32.20200824.3.0-live-kernel-x86_64",
		},
		"openstack release": {
			run: func(s Stream) string {
				return s.Architectures["x86_64"].Artifacts["metal"].Release
			},
			want: "32.20200824.3.0",
		},
		"aws image": {
			run: func(s Stream) string {
				return s.Architectures["x86_64"].Images.AWS.Regions["us-west-2"].Image
			},
			want: "ami-002ac5b87eb32f650",
		},
		"gcp image": {
			run: func(s Stream) string {
				return s.Architectures["x86_64"].Images.GCP.Name
			},
			want: "fedora-coreos-32-20200824-3-0-gcp-x86-64",
		},
	}

	for name, c := range cases {
		have := c.run(*strm)
		if c.want != have {
			t.Errorf("Resolve case: %v; have %v; want %v", name, have, c.want)
		}
	}

}
