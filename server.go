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

	// CancelArgs is used to cancel a notification
	CancelArgs struct {}

	// A Result is returned back to the client.
	Result bool

	Server struct {
		Remotes map[string]chan bool
		NotificationCancel chan bool
	}
)

// NewServer creates a new server object.
func NewServer(remotes map[string]chan bool, cancelNotification chan bool) *Server {
	return &Server{Remotes: remotes, NotificationCancel: cancelNotification}
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


// CancelNotification is used to cancel the running notification
func (s *Server) CancelNotification(args CancelArgs, result *Result) error {
	s.NotificationCancel <- true	
	*result = true
	return nil
}
