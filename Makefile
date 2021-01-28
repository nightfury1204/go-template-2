build:
	./build.sh

build-proto:
	./build-proto.sh

run: build
	go run main.go serve-grpc-rest --config app.config.yaml

serve:
	docker-compose down
	docker-compose up -d
