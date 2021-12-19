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

type HandlerFunc func(conn net.Conn)
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
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != io.EOF {
		log.Printf("%s", buf[:n])
	}

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

	_, err = url.ParseRequestURI(parts[1])
	if err != nil {
		log.Print("url ParseRequestURI err: ", err)
	}
	s.mu.RLock()
	if handler, ok := s.handlers[parts[1]]; ok {
		s.mu.RUnlock()
		handler(conn)
	}
}
