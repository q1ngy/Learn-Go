package main

import (
	"context"
	"hello_server/pb"
	"io"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloResp, error) {
	reply := "hello " + in.Name
	return &pb.HelloResp{Reply: reply}, nil
}

// LotsOfReplies 返回使用多种语言打招呼
func (s *server) LotsOfReplies(in *pb.HelloReq, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}

	for _, word := range words {
		data := &pb.HelloResp{
			Reply: word + in.GetName(),
		}
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) LotsOfGreetings(stream grpc.ClientStreamingServer[pb.HelloReq, pb.HelloResp]) error {
	// pb.Greeter_LotsOfGreetingsServer
	reply := "Hello: "
	for {
		// 接收客户端发来的流式数据
		res, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			return stream.SendAndClose(&pb.HelloResp{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += res.GetName() + " "
	}
}

// BidiHello 双向流式打招呼
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := magic(in.GetName()) // 对收到的数据做些处理

		// 返回流式响应
		if err := stream.Send(&pb.HelloResp{Reply: reply}); err != nil {
			return err
		}
	}
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
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
