package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	Conn               *grpc.ClientConn
	SpreadsheetService proto.SpreadsheetServiceClient
}

func NewGRPCClients() *GRPCClients {
	grpcURL := os.Getenv("GRPC_SERVER_URL")
	if grpcURL == "" {
		log.Fatal("GRPC_SERVER_URL not found in environment variables")
	}

	conn, err := grpc.NewClient(
		grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// ⭐ CRITICAL: Dynamically balances requests across expanding K8s replicas
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)

	if err != nil {
		log.Fatalf("Failed to initialize gRPC client connection: %v", err)
	}

	return &GRPCClients{
		Conn:               conn,
		SpreadsheetService: proto.NewSpreadsheetServiceClient(conn),
	}
}
