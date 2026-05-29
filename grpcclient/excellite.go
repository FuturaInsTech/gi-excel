package grpcclient

import (
	"log"
	"os"

	"github.com/FuturaInsTech/gi-excel/proto"
	"google.golang.org/grpc"
)

func NewExcelLiteGRPCClient() proto.SpreadsheetServiceClient {

	grpcURL := os.Getenv("EXCEL_LITE_GRPC_SERVER_URL")

	if grpcURL == "" {

		log.Fatal(
			"EXCEL_LITE_GRPC_SERVER_URL not found",
		)
	}

	conn, err := grpc.Dial(
		grpcURL,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {

		log.Fatalf(
			"Failed to connect to Excel Lite gRPC: %v",
			err,
		)
	}

	return proto.NewSpreadsheetServiceClient(conn)
}
