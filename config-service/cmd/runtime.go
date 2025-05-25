package cmd

import (
	"github.com/spf13/cobra"
	grpcservercmd "github.com/tuandq2112/go-microservices/config-service/cmd/grpc_server"
)

var RootCmd = &cobra.Command{
	Use:   "config-service",
	Short: "Config service",
	Long:  "Config service",
}


func init() {
	RootCmd.AddCommand(grpcservercmd.GrpcServerCmd)
}