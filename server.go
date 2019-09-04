package main

import (
	"errors"
	"fmt"
)

type (
	// Args is passed during rpc. It is not currently
	// used.
	Args struct {
		Target string
	}

	// A Result is returned back to the client.
	Result bool

	Server struct {
		Remotes map[string]chan bool
	}
)

// NewServer creates a new server object.
func NewServer(remotes map[string]chan bool) *Server {
	return &Server{Remotes: remotes}
}

// Call is used to remotely trigger an update
func (s *Server) Call(args Args, result *Result) error {
	remote, ok := s.Remotes[args.Target]
	if !ok {
		*result = false
		return errors.New("No remote with that name")
	}
	fmt.Printf("Triggering remote: %s\n", args.Target)
	remote <- true
	*result = true
	return nil
}
