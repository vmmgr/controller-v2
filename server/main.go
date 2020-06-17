package server

// Client gRPC Server

import (
	pb1 "github.com/vmmgr/controller/proto/proto-go"
	pb2 "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

const basePort = ":50200"
const vmPort = ":50210"

type baseServer struct {
	pb1.UnimplementedControllerServer
}
type vmServer struct {
	pb2.UnimplementedNodeServer
}

func BaseServer() {
	lis, err := net.Listen("tcp", basePort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb1.RegisterControllerServer(s, &baseServer{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

func VMServer() {
	lis, err := net.Listen("tcp", vmPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb2.RegisterNodeServer(s, &vmServer{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
