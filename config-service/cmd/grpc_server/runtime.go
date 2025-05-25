package grpcserver

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tuandq2112/go-microservices/config-service/appconfig"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "Start the gRPC server",
	Long:  "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()

		server, err := InitializeServer()
		if err != nil {
			log.Fatalf("Failed to initialize server: %v", err)
		}
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}
