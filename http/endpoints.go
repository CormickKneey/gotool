package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Endpoint ...
type Endpoint struct {
	raw string
	url *url.URL
}

// Apply modify the request
func (ep *Endpoint) Apply(r *http.Request) {
	r.URL.Host = ep.url.Host

	// r.Host = ep.url.Host

	// TODO: anything else?
}

// NewEndpoint ...
func NewEndpoint(raw string) (ep *Endpoint, err error) {

	ep = &Endpoint{
		raw: raw,
	}

	ep.url, err = url.Parse(raw)

	return
}

// HealthChecker ... TODO: split to interface
type HealthChecker func(ep *Endpoint, req *http.Request, resp *http.Response) bool

// NewSimpleHealthCheck ... TODO: how to add timeout
func NewSimpleHealthCheck(method string, path string, content string) HealthChecker {
	return func(ep *Endpoint, req *http.Request, resp *http.Response) bool {
		
		if resp == nil {
			// modify the req
			req.URL.Path = path
			req.Method = method

			// just prepare, return nothings
			return true
		}

		defer resp.Body.Close()

		// must check the reult
		if content != "" {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return false
			}

			return bytes.Equal(data, []byte("ok"))
		}

		// let check the http status code
		return resp.StatusCode > http.StatusOK && resp.StatusCode < http.StatusPartialContent
	}
}