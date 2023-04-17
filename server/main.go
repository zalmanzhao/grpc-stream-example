package main

import (
	"google.golang.org/grpc"
	"grpc.server/controllers"
	"grpc.server/protos"
	"grpc.server/storage"
	"log"
	"net"
)

func main() {
	address := ":" + "8888"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println(err)
	}
	server := grpc.NewServer()

	// Bootstrap upload server.
	uplSrv := controllers.NewUploadServer(storage.New("tmp/"))
	// Register and start gRPC server.
	protos.RegisterUploadServer(server, &uplSrv)
	protos.RegisterMessageServer(server, &controllers.Message{})
	err = server.Serve(listener)
	if err != nil {
		log.Println(err)
	}
}
