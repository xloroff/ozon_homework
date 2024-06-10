package grpcserver

import "google.golang.org/grpc"

// APIHandler интерфейс реализующий добавление ручен.
type APIHandler interface {
	RegisterGrpcServer(server *grpc.Server)
}

// RegisterAPI добавляет ручки.
func (s *server) RegisterAPI(h []APIHandler) {
	for _, handl := range h {
		handl.RegisterGrpcServer(s.grpcServer)
	}
}
