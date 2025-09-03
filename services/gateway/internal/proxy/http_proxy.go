package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(targetBase string, stripPrefix string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(targetBase)
	if err != nil {
		return nil, err
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	origDirector := rp.Director
	rp.Director = func(r *http.Request) {
		origDirector(r)
		if stripPrefix != "" && strings.HasPrefix(r.URL.Path, stripPrefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, stripPrefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}
	}
	return rp, nil
}
