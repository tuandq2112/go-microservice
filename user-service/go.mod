module github.com/tuandq2112/go-microservices/user-service

go 1.23.0

require (
	github.com/oklog/run v1.1.0
	github.com/tuandq2112/go-microservices/shared v0.0.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.72.1
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/tuandq2112/go-microservices/shared => ../shared
