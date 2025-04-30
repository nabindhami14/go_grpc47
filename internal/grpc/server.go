package grpc

import newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"

type Server struct {
	newsv1.UnimplementedNewsServiceServer
}

func NewServer() *Server {
	return &Server{}
}
