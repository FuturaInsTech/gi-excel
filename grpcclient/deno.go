package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/denoproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDenoGRPCClient() (denoproto.FunctionRuntimeClient, *grpc.ClientConn) {

	grpcURL := os.Getenv("DENO_GRPC_SERVER_URL")

	if grpcURL == "" {

		log.Fatal(
			"DENO_GRPC_SERVER_URL not found",
		)
	}

	conn, err := grpc.NewClient(
		grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// ⭐ CRITICAL: Dynamically balances requests across expanding K8s replicas
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)

	if err != nil {

		log.Fatalf(
			"Failed to initialize connection to Deno Runtime gRPC: %v",
			err,
		)
	}

	return denoproto.NewFunctionRuntimeClient(conn), conn
}
