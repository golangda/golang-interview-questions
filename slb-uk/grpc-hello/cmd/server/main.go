package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/slb-uk/grpc-hello/api/hellopb"
)

type greeterServer struct {
	hellopb.UnimplementedGreeterServer
}

// Unary RPC
func (g *greeterServer) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	name := req.GetName()
	return &hellopb.HelloResponse{Message: fmt.Sprintf("Hello, %s! 👋", name)}, nil
}

// Server-streaming RPC
func (g *greeterServer) GreetManyTimes(req *hellopb.HelloRequest, stream hellopb.Greeter_GreetManyTimesServer) error {
	name := req.GetName()
	for i := 1; i <= 5; i++ {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
		}
		msg := fmt.Sprintf("[%d/5] Hello, %s!", i, name)
		if err := stream.Send(&hellopb.HelloResponse{Message: msg}); err != nil {
			return err
		}
		time.Sleep(600 * time.Millisecond)
	}
	return nil
}

// --- Interceptors ---
func unaryLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[UNARY] method=%s dur=%s err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func authUnaryInterceptor(validToken string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if validToken == "" { // auth disabled
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}
		vals := md.Get("authorization")
		expected := "Bearer " + validToken
		if len(vals) == 0 || vals[0] != expected {
			return nil, fmt.Errorf("unauthorized")
		}
		return handler(ctx, req)
	}
}

func main() {
	addr := ":50051"
	if v := os.Getenv("GRPC_ADDR"); v != "" {
		addr = v
	}
	token := os.Getenv("GREETER_TOKEN") // optional

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryLoggerInterceptor,
			authUnaryInterceptor(token),
		),
	)

	hellopb.RegisterGreeterServer(s, &greeterServer{})

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("shutting down gracefully...")
		s.GracefulStop()
	}()

	log.Printf("gRPC server listening on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}