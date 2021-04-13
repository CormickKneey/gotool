package http

import (
	"testing"
)

func TestServer_Run(t *testing.T) {
	fakeServer := NewServer(":8081")
	fakeServer.Run()

}
