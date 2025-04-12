package httpserver

import (
	"net"
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func Prefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

func ReadTimeout(time time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = time
	}
}

func WriteTimeout(time time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = time
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
