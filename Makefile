.PHONY: build
build:
	go build -o ./bin/server ./cmd/server
	go build -o ./bin/client ./cmd/client

.PHONY: gen-proto
gen-proto:
	@ rm -rf gen/stt
	@ protoc --proto_path=proto --go_out=gen \
	 --go_opt=paths=source_relative \
	 --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	  proto/stt/service/v1/stt.proto

.PHONY: build-stt-docker
build-stt-docker:
	@ docker build ./cmd/stt/Dockerfile -t dss:latest

.PHONY: gen-db
gen-db:
	@ go generate ./gen/ent

# Replace local assets with GitHub one
.PHONY: run-stt-docker
run-stt-docker:
	docker run --rm -it --name ds \
	-v /home/tas/Documents/GoLang/YPracticum/coqui-ai-assets:/opt/deepspeech \
	-p 0.0.0.0:9090:9090 dss:latest

# Generate mocks
MOCKS_DESTINATION=mock
.PHONY: mocks
mocks: internal/db/*.go
	@echo "Generating mocks for db package..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/db_mocks.go -package=mock; done