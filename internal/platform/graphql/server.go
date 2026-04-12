package graphql

import (
	"fmt"
)

type Server struct {
	port                int
	enablePlayground    bool
	enableIntrospection bool
}

func New(port int) *Server {
	return &Server{
		port:                port,
		enablePlayground:    true,
		enableIntrospection: true,
	}
}

func (s *Server) EnablePlayground(enabled bool) {
	s.enablePlayground = enabled
}

func (s *Server) EnableIntrospection(enabled bool) {
	s.enableIntrospection = enabled
}

func (s *Server) Address() string {
	return fmt.Sprintf(":%d", s.port)
}
