package http

import (
	"context"
	"io"
	"log"
	"net/http"
)

// SimpleServer a simple http server to test the RoundTrip
type SimpleServer struct {
	// use it with default mutex
	http.Server
	addition interface{}

	shutdownCh chan struct{}
}

func NewServer(addr string) *SimpleServer{
	return &SimpleServer{
		Server:   http.Server{
			Addr:              addr,
			Handler:           nil,
		},
		addition: nil,
		shutdownCh: make(chan struct{},1),
	}
}


func (s *SimpleServer) Run() {
	// use default mutex
	http.HandleFunc("/api", HelloServer)
	http.HandleFunc("/shutdown",s.HandleShutdown)

	go func() {
		err := s.Server.ListenAndServe()
		if err != nil {
			log.Println("start server error ",err)
			return
		}
	}()

	<-s.shutdownCh
	s.ShutDown()
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request: %+v",req)
	io.WriteString(w, "hello, world!\n")
}

func (s *SimpleServer)HandleShutdown(w http.ResponseWriter, req *http.Request){
	s.shutdownCh <- struct{}{}
}

func (s *SimpleServer)ShutDown(){
	err := s.Server.Shutdown(context.Background())
	if err != nil {
		log.Println("shut down server error ",err)
		return
	}
}

