package server

import "context"

type Server struct {
}

func NewServer() (*Server, error) {

	return &Server{}, nil
}

func (srv *Server) Run(ctx context.Context) error {

	return nil
}
