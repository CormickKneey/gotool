package http

import (
	"testing"
)

const url1 = "localhost:9880"
const url2 = "localhost:9881"
const url3 = "localhost:9882"

func TestServer_Run(t *testing.T) {
	stopCh := make(chan struct{})
	fakeServer1 := NewServer(url1)
	fakeServer2 := NewServer(url2)
	fakeServer3 := NewServer(url3)
	go fakeServer1.Run()
	go fakeServer2.Run()
	go fakeServer3.Run()
	<-stopCh
}
