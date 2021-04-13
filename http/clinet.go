package http

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type SimpleClient struct {
	http.Client
}
type Balancer struct {
	http.RoundTripper
	hosts[] string

	sync.Mutex
}

func NewBalance(hosts []string) *Balancer {
	return &Balancer{
		RoundTripper: http.DefaultTransport,
		hosts:        hosts,
		Mutex:        sync.Mutex{},
	}
}

func (b *Balancer)RoundTrip(req *http.Request) (*http.Response, error){
	var req2 *http.Request
	// copy request
	*req2 = *req
	resp,err := b.RoundTripper.RoundTrip(req)
	if err != nil {
		log.Println("get resp err",err)

	}
	return resp,nil
}

func (b *Balancer)GetHealthyOne() string{
	randIndex := rand.Perm(len(b.hosts))
	for _,index := range randIndex {
		resp,err := http.Get(b.hosts[index]+"/health_z")
		if err != nil {
			log.Println("get response error ",err)
			continue
		}
		data,_ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if string(data) == "ok" {
			return b.hosts[index]
		}
	}
	return b.hosts[0]
}