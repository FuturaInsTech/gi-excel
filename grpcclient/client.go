package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/proto"
	"google.golang.org/grpc"
)

type GRPCClients struct {
	/*Conn               *grpc.ClientConn*/
	SpreadsheetService proto.SpreadsheetServiceClient
}

func NewGRPCClients() *GRPCClients {
	grpcURL := os.Getenv("GRPC_SERVER_URL")
	if grpcURL == "" {
		log.Fatal("GRPC_SERVER_URL not found in environment variables")
	}

	conn, err := grpc.Dial(grpcURL,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC: %v", err)
	}

	return &GRPCClients{
		/*Conn:               conn,*/
		SpreadsheetService: proto.NewSpreadsheetServiceClient(conn),
	}
}
