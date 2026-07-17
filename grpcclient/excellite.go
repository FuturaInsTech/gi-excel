package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExcelLiteGRPCClient() (proto.SpreadsheetServiceClient, *grpc.ClientConn) {

	grpcURL := os.Getenv("EXCEL_LITE_GRPC_SERVER_URL")

	if grpcURL == "" {

		log.Fatal(
			"EXCEL_LITE_GRPC_SERVER_URL not found",
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
			"Failed to initialize connection to Excel Lite gRPC: %v",
			err,
		)
	}

	return proto.NewSpreadsheetServiceClient(conn), conn
}
