package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedGreeterServer
	mu    sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloResp, error) {
	// metadata
	defer func() {
		// 发送结束后发送trailer
		trailer := metadata.Pairs(
			"timestamp", strconv.Itoa(int(time.Now().Unix())),
		)
		grpc.SetTrailer(ctx, trailer)
	}()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "无效请求")
	}
	v := md.Get("token")
	if len(v) < 1 || v[0] != "my-token" {
		return nil, status.Error(codes.Unauthenticated, "无效 token")
	}
	//if v, ok := md["token"]; ok {

	name := in.GetName()
	s.mu.Lock()
	s.count[name]++
	s.mu.Unlock()
	if s.count[name] > 1 {
		// grpc status
		st := status.New(codes.ResourceExhausted, "name request limit")
		// detail
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "限制每个name调用一次",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}

		return nil, ds.Err()
	}

	reply := "hello " + name

	// 发送数据前发送header
	header := metadata.New(map[string]string{
		"location": "北京",
	})
	grpc.SetHeader(ctx, header)

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
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	err = s.Serve(l)
	if err != nil {
		log.Printf("failed to server, err: %v\n", err)
		return
	}
}
