package rungrpc

import (
	"log"

	"github.com/spf13/cobra"
)

var RunGRPCServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "Run GRPC Server",
	Long:  "Run GRPC Server",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := InitializeApp()
		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}
		app.Start()
	},
}
