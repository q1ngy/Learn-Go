package main

import (
	"context"
	"gateway/server/pb"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{Reply: "hello " + in.Name}, nil
}

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("failed to listen, err: %v\n", err)
		return
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterGreeterServer(s, &server{})
	//err = s.Serve(l) // 阻塞，就执行不到下边的代码了
	//if err != nil {
	//	log.Printf("failed to server, err: %v\n", err)
	//	return
	//}

	go func() {
		log.Fatalln(s.Serve(l))
	}()

	conn, err := grpc.NewClient(
		"127.0.0.1:9000",
		grpc.WithBlock(), // 阻塞直到连接成功
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// 注册Greeter
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	// 8090端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
