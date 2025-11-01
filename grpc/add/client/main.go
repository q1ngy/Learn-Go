package main

import (
	"add_client/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("new client err: %v\n", err)
	}
	defer conn.Close()

	client := proto.NewCalcServiceClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	res, err := client.Add(ctx, &proto.AddReq{X: 1, Y: 2})
	if err != nil {
		log.Fatalf("client add err: %v\n", err)
	}
	fmt.Println("result is ", res.GetResult())
}
