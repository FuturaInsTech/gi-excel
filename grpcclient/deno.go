package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/denoproto"
	"google.golang.org/grpc"
)

func NewDenoGRPCClient() denoproto.FunctionRuntimeClient {

	grpcURL := os.Getenv("DENO_GRPC_SERVER_URL")

	if grpcURL == "" {

		log.Fatal(
			"DENO_GRPC_SERVER_URL not found",
		)
	}

	conn, err := grpc.Dial(
		grpcURL,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {

		log.Fatalf(
			"Failed to connect to Deno Runtime gRPC: %v",
			err,
		)
	}

	return denoproto.NewFunctionRuntimeClient(conn)
}
