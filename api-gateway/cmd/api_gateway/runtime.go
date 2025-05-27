package api_gateway

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/go-microservices/api-gateway/appconfig"
	"github.com/tuandq2112/go-microservices/api-gateway/infrastructure/server"
)

var HttpGatewayCmd = &cobra.Command{
	Use:   "http-gateway",
	Short: "HTTP Gateway",
	Long:  "HTTP Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()

		httpServer := server.NewHttpServer()
		httpServer.Start()
	},
}
