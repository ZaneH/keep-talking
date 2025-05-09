gen-protobuf:
	@echo "Generating protobuf files..."
	@protoc \
		--go_out=./internal/infrastructure/grpc \
  	--go_opt=paths=source_relative \
		--go-grpc_out=./internal/infrastructure/grpc \
  	--go-grpc_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
		proto/*proto
	@echo "Protobuf files generated successfully."

test:
	@go test ./... -v