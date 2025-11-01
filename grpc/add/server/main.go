package main

import (
	"add_server/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedCalcServiceServer
}

func (s *server) Add(ctx context.Context, in *proto.AddReq) (*proto.AddResp, error) {
	res := in.GetX() + in.GetY()
	return &proto.AddResp{Result: res}, nil
}

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("net listen err: %v\n", err)
	}
	s := grpc.NewServer()
	proto.RegisterCalcServiceServer(s, &server{})
	err = s.Serve(l)
	if err != nil {
		log.Fatalf("server err: %v\n", err)
	}
}
