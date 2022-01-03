package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

type Request struct {
	Conn        net.Conn
	QueryParams url.Values
}

type HandlerFunc func(req *Request)
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	//start server on the given address addr

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Print(err)
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != io.EOF {
		log.Printf("%s", buf[:n])
	}

	var req Request
	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
		log.Print("requestLineEndErr: ", requestLineEnd)
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		log.Print("partsErr: ", parts)
	}

	_, path, version := parts[0], parts[1], parts[2]
	if version != "HTTP/1.1" {
		log.Print("dsf")
	}
	uri, _ := url.ParseRequestURI(path)
	log.Print(uri.Path)
	log.Print(uri.Query())
	_, err = url.ParseRequestURI(parts[1])
	if err != nil {
		log.Print("url ParseRequestURI err: ", err)
	}

	req.Conn = conn
	req.QueryParams = uri.Query()

	s.mu.RLock()
	if handler, ok := s.handlers[parts[1]]; ok {
		s.mu.RUnlock()
		handler(&req)
	}
}
