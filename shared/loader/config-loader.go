package appconfig

import (
	"context"
	"fmt"
	"log"

	configPb "github.com/tuandq2112/go-microservices/shared/proto/types/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func LoadConfigFromConfigService(configServiceHost string, configServicePort string, serviceName string, env string) (map[string]any, error) {
	ctx := context.Background()
	target := fmt.Sprintf("%s:%s", configServiceHost, configServicePort)
	conn, err := grpc.NewClient( target, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 10)),
	)		
	if err != nil {
		log.Fatalf("Failed to connect to config service: %v", err)
	}
	defer conn.Close()

	client := configPb.NewConfigServiceClient(conn)
	
	clientConfig, err := client.GetConfig(ctx, &configPb.GetConfigRequest{
		ServiceName: serviceName,
		Env:         env,
	})
	if err != nil {
		return nil, err
	}

	return clientConfig.Value.AsMap(), nil
}