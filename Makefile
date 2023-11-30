.PHONY: build
build:
	go build -o ./bin/server ./cmd/server
	go build -o ./bin/client ./cmd/client

.PHONY: gen-proto
gen-proto:
	protoc --proto_path=proto --go_out=gen \
	 --go_opt=paths=source_relative \
	 --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	  proto/stt/service/v1/service.proto