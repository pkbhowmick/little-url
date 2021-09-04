package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	pb "github.com/pkbhowmick/url-lite/grpc/url_gen/urlgen"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50001"
)

type gRPCServer struct {
	pb.UnimplementedKeyGenServer
}

func (s *gRPCServer)GenerateKey(ctx context.Context, empty *empty.Empty) (*pb.Key, error) {
	key := uuid.New().String()
	key = key[:8]
	return &pb.Key{Key: key}, nil
}

func Start() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterKeyGenServer(s, &gRPCServer{})
	log.Printf("grpc server is listening at %v",lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}