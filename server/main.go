package main

import (
	"google.golang.org/grpc"
	"grpc.server/controllers"
	"grpc.server/protos"
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
	protos.RegisterMessageServer(server, &controllers.Message{})
	err = server.Serve(listener)
	if err != nil {
		log.Println(err)
	}
}
