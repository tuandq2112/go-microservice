package cmd

import (
	"github.com/spf13/cobra"
	api_gateway "github.com/tuandq2112/go-microservices/api-gateway/cmd/api_gateway"
)

var RootCmd = &cobra.Command{
	Use: "api-gateway",
	Short: "API Gateway",
	Long: "API Gateway",
}

func init() {
	RootCmd.AddCommand(api_gateway.HttpGatewayCmd)
}