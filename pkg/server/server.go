package server

import (
	"context"
	"device-info-sorting-storage/pkg/middleware"
)

type Server struct {
}

func NewServer() (*Server, error) {

	return &Server{}, nil
}

func (srv *Server) Run(ctx context.Context) error {
	return middleware.HandleDeviceInfo(ctx)
}
