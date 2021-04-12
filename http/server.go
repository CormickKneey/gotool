package http

import (
	"io"
	"log"
	"net/http"
)

// Server a simple http server
type Server struct {
	Addr string `json:"addr"`
	MaxCoon int `json:"max_coon"`
	mux *http.ServeMux
}

func NewServer() *Server{
	return &Server{
		Addr:    ":8081",
		MaxCoon: 10,
		mux:     http.NewServeMux(),
	}
}

// wrapper the mux.handler
func (server *Server) httpHandler(pattern string, handler func(http.ResponseWriter, *http.Request))  {
	server.mux.HandleFunc(pattern,handler)
}

func (server *Server) Run() {
	server.httpHandler("/api", HelloServer)

	err := http.ListenAndServe(server.Addr, server.mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request: %+v",req)
	io.WriteString(w, "hello, world!\n")
}

