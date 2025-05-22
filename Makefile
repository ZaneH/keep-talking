gen-protobuf:
	@echo "Generating protobuf files..."
	@buf generate
	@echo "Protobuf files generated successfully."

swagger-ui:
	@docker run -p 80:8080 \
    -e SWAGGER_JSON=/proto/game.swagger.json \
    -v $(PWD)/openapiv2/proto:/proto \
    swaggerapi/swagger-ui
	@echo "Swagger UI is running at http://localhost:80"

test:
	@go test ./... -v