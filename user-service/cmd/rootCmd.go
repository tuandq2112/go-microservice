package cmd

import (
	"github.com/spf13/cobra"
	rungrpc "github.com/tuandq2112/go-microservices/user-service/cmd/run_grpc"
)

var RootCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "API Gateway",
	Long:  "API Gateway",
}

func init() {
	RootCmd.AddCommand(rungrpc.RunGRPCServerCmd)
}
