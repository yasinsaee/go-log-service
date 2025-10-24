package app

import (
	"log"
	"net"

	"github.com/yasinsaee/go-log-service/internal/app/config"
	loggrpc "github.com/yasinsaee/go-log-service/internal/handlers/grpc/log"
	logpb "github.com/yasinsaee/go-log-service/log-service/log"
	"github.com/yasinsaee/go-log-service/pkg/elastic"
	"google.golang.org/grpc"
)

func StartGRPCServer() {
	port := config.GetEnv("PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// //handlers
	logHandler := loggrpc.New(elastic.ClientInstance)

	// //register grpc services
	logpb.RegisterLogServiceServer(s, logHandler)

	log.Printf("gRPC server is running on port %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
