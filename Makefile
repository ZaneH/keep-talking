gen-protobuf:
	@echo "Generating protobuf files..."
	@protoc \
		--go_out=./internal/infrastructure/grpc \
  	--go_opt=paths=source_relative \
		--go-grpc_out=./internal/infrastructure/grpc \
  	--go-grpc_opt=paths=source_relative \
		proto/*proto
	@echo "Protobuf files generated successfully."