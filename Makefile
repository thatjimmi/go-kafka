websocket:
	go run ./cmd/websocket .
producer:
	go run ./cmd/producer .
consumer:
	go run ./cmd/consumer .
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative ./proto/*.proto

