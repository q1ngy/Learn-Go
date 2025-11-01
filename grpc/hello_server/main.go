package main

import (
	"context"
	"hello_server/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloResp, error) {
	reply := "hello " + in.Name
	return &pb.HelloResp{Reply: reply}, nil
}

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("failed to listen, err: %v\n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(l)
	if err != nil {
		log.Printf("failed to server, err: %v\n", err)
		return
	}
}
