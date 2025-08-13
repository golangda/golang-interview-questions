package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/slb-uk/grpc-hello/api/hellopb"
)

func main() {
	addr := "localhost:50051"
	if v := os.Getenv("GRPC_ADDR"); v != "" {
		addr = v
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := hellopb.NewGreeterClient(conn)

	// Prepare metadata (auth token optional)
	var ctx context.Context = context.Background()
	if tok := os.Getenv("GREETER_TOKEN"); tok != "" {
		md := metadata.New(map[string]string{"authorization": "Bearer " + tok})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Unary with timeout
	uctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := client.SayHello(uctx, &hellopb.HelloRequest{Name: "Rahul"})
	if err != nil {
		log.Fatalf("SayHello: %v", err)
	}
	fmt.Println("Unary:", res.GetMessage())

	// Server-streaming
	stream, err := client.GreetManyTimes(ctx, &hellopb.HelloRequest{Name: "Rahul"})
	if err != nil {
		log.Fatalf("GreetManyTimes: %v", err)
	}
	fmt.Println("Stream:")
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream recv: %v", err)
		}
		fmt.Println(" ", msg.GetMessage())
	}
}