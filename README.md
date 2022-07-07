# grpc-stream-example
grpc stream example code


go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=./server --go-grpc_out=./server ./server/proto/*.proto
protoc --go_out=./client --go-grpc_out=./client ./client/proto/*.proto
