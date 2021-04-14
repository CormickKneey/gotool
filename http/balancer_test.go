package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const endpoint1 = "http://localhost:9880"
const endpoint2 = "http://localhost:9881"
const endpoint3 = "http://localhost:9882"

type simpleClient struct {
	http.Client
}

func (s *simpleClient) Do(url string) {
	resp, err := s.Get(url + "/name")
	if err != nil {
		log.Println("get response error ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response error ", err)
	}
	log.Println("Response: ", string(body))
}

func TestNewBalancerRoundTripper(t *testing.T) {
	brt := NewBalancerRoundTripper(endpoint1, endpoint2, endpoint3)
	cli := &simpleClient{
		Client: http.Client{Transport: brt(http.DefaultTransport)},
	}
	for true {
		cli.Do(endpoint1)
	}
}
