package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/thatjimmi/go-kafka/proto/services"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExampleServiceServer
}

func startGRPCServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, &server{})
	log.Printf("starting gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(grpcAddress, httpAddress string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterExampleServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	log.Printf("starting HTTP server on %s", httpAddress)
	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}

func main() {
	go startGRPCServer(":8080")                // Start gRPC server on port 8080
	startHTTPServer("localhost:8080", ":8081") // Start HTTP server on port 8081
}
