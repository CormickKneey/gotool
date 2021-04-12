package http

import (
	"testing"
)

func TestServer_Run(t *testing.T) {
	fakeServer := NewServer()
	fakeServer.Run()
}
