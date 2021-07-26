protoc --go_out=./gen --go-grpc_out=./gen --proto_path=./protobuf mymsg.proto
go mod tidy
go mod vendor