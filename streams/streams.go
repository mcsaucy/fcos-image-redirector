package streams

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	streamCachePrefix = "stream cache for: "
)

type Arch struct {
	Artifacts map[string]Artifact `json:"artifacts"`
	Images    CloudImages         `json:"images"`
}

type CloudImages struct {
	AWS AWSImages `json:"aws"`
	GCP GCPImages `json:"gcp"`
}

type AWSRegionalOffering struct {
	Release string `json:"release"`
	Image   string `json:"image"`
}

type AWSImages struct {
	Regions map[string]AWSRegionalOffering `json:"regions"`
}

type GCPImages struct {
	Project string `json:"project"`
	Family  string `json:"family"`
	Name    string `json:"name"`
}

type StreamMetadata struct {
	LastModified string `json:"last-modified"`
}

type Formats map[string]map[string]*Resource

type Artifact struct {
	Release string `json:"release"`
	Formats `json:"formats"`
}

type Resource struct {
	Location  string `json:"location"`
	Signature string `json:"signature"`
	Sha256    string `json:"sha256"`
}

type Stream struct {
	Name          string          `json:"stream"`
	Metadata      StreamMetadata  `json:"metadata"`
	Architectures map[string]Arch `json:"architectures"`
}

type Resolver struct {
	cli   *http.Client
	cache *cache.Cache
}

func New() Resolver {
	return Resolver{
		cli:   &http.Client{Timeout: 30 * time.Second},
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (r Resolver) Resolve(ctx context.Context, stream string) (*Stream, error) {
	key := streamCachePrefix + stream
	cached, found := r.cache.Get(key)
	if found {
		return cached.(*Stream), nil
	}

	s, err := r.resolve(ctx, stream)
	if err != nil {
		return nil, err
	}
	r.cache.Set(key, s, cache.DefaultExpiration)
	return s, nil
}

func (r Resolver) resolve(ctx context.Context, stream string) (*Stream, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "builds.coreos.fedoraproject.org",
		Path:   "/streams/" + stream + ".json",
	}
	req := http.Request{
		Method: "GET",
		URL:    u,
	}
	req = *req.WithContext(ctx)

	resp, err := r.cli.Do(&req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %v: %w", u, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-OK status when fetching %v: %v", u, resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of %v: %w", u, err)
	}

	strms := &Stream{}
	err = json.Unmarshal(bodyBytes, strms)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal contents of %v: %w", u, err)
	}

	return strms, nil
}
