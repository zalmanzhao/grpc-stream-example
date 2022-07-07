package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc.server/controllers"
	"grpc.server/protos"
	"net"
)

func main() {
	address := ":" + "8888"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	server := grpc.NewServer()
	protos.RegisterMessageServer(server, &controllers.Message{})
	err = server.Serve(listener)
	if err != nil {
		fmt.Println(err)
	}
}
