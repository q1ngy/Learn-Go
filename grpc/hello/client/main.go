package main

import (
	"bufio"
	"context"
	"fmt"
	"hello_client/pb"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gprc new client failed, err: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	//sayHello(c, ctx)
	//streamSayHello(c, ctx)
	//clientStream(c, ctx)
	//runBidiHello(c, ctx)
	sayHelloWithMD(c, ctx)
}

func clientStream(c pb.GreeterClient, ctx context.Context) {
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		fmt.Printf("stream err: %v\n", err)
	}
	names := []string{"slim", "nolan", "qingy"}
	for _, name := range names {
		stream.Send(&pb.HelloReq{Name: name})
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("stream recv err: %v\n", err)
	}
	fmt.Printf("stream result: %v\n", recv)
}

func sayHello(c pb.GreeterClient, ctx context.Context) {
	resp, err := c.SayHello(ctx, &pb.HelloReq{Name: "slim"})
	if err != nil {
		log.Fatalf("SayHello failed, err: %v\n", err)
	}
	log.Printf("resp: %v", resp.GetReply())
}

func streamSayHello(c pb.GreeterClient, ctx context.Context) {
	stream, err := c.LotsOfReplies(ctx, &pb.HelloReq{Name: "slim"})
	if err != nil {
		log.Fatalf("hello stream failed, err: %v\n", err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("stream recv err: %v\n", err)
			return
		}
		log.Printf("stream res: %v\n", res.GetReply())
	}
}

func runBidiHello(c pb.GreeterClient, ctx context.Context) {
	// 双向流模式
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed, err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI：%s\n", in.GetReply())
		}
	}()
	// 从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') // 读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&pb.HelloReq{Name: cmd}); err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func sayHelloWithMD(c pb.GreeterClient, ctx context.Context) {
	md := metadata.Pairs(
		"token", "my-token",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header, trailer metadata.MD
	resp, err := c.SayHello(ctx, &pb.HelloReq{Name: "slim"},
		grpc.Header(&header),
		grpc.Trailer(&trailer))
	fmt.Printf("metadata header: %v, trailer: %v\n", header, trailer)

	if err != nil {
		log.Fatalf("SayHello failed, err: %v\n", err)
	}
	log.Printf("resp: %v", resp.GetReply())

}
