package http

import (
	"context"
	"io"
	"log"
	"net/http"
)

// SimpleServer a simple http server to test the RoundTrip
type SimpleServer struct {
	http.Server
	addition interface{}

	shutdownCh chan struct{}
	mux *http.ServeMux
}

func NewServer(addr string) *SimpleServer{
	mux := http.NewServeMux()
	return &SimpleServer{
		Server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		addition:   nil,
		shutdownCh: make(chan struct{}, 1),
		mux:        mux,
	}
}


func (s *SimpleServer) Run() {
	// use default mutex ,  noop!

	s.mux.HandleFunc("/health_z", HelloServer)
	s.mux.HandleFunc("/shutdown",s.HandleShutdown)
	s.mux.HandleFunc("/name",s.HandleName)

	go func() {
		log.Println("Started server on ",s.Addr)
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
	io.WriteString(w, "ok")
}

func (s *SimpleServer)HandleShutdown(w http.ResponseWriter, req *http.Request){
	s.shutdownCh <- struct{}{}
}

func (s *SimpleServer)HandleName(w http.ResponseWriter, req *http.Request){
	log.Println("[request] ",req)
	io.WriteString(w, "Here is " + s.Addr)
}

func (s *SimpleServer)ShutDown(){
	err := s.Server.Shutdown(context.Background())
	if err != nil {
		log.Println("shut down server error ",err)
		return
	}
}

