websocket:
	go run ./cmd/websocket .
producer:
	go run ./cmd/producer .
consumer:
	go run ./cmd/consumer .
grpc:
	go run ./cmd/grpc .
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative ./proto/*.proto
test:
	go test -v ./...