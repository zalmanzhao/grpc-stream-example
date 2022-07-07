package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc.client/protos"
	"io"
	"log"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	c := protos.NewMessageClient(conn)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	send(c)
	sendServerStream(c)
	sendClientStream(c)
	sendClientServerStream(c)
}
func send(c protos.MessageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Send(ctx, &protos.SendMessage{Name: "UnaryAPI"})
	if err != nil {
		panic(err)
	}
	log.Printf("Revice Ms: %s", r)
}
func sendServerStream(c protos.MessageClient) {
	stream, err := c.SendServerStream(context.Background(), &protos.SendMessage{Name: "ServerStream"})
	if err != nil {
		panic(err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Recv error:%v", err)
			continue
		}
		log.Printf("Recv Ms:%v", resp)
	}
}
func sendClientStream(c protos.MessageClient) {
	stream, err := c.SendClientStream(context.Background())
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		err := stream.Send(&protos.SendMessage{Name: "ClientStream" + strconv.Itoa(i)})
		if err != nil {
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("failed to recv: %v", err)
	}
	log.Printf("Recv Ms: %s", resp)
}
func sendClientServerStream(c protos.MessageClient) {
	stream, err := c.SendClientServerStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("recv-server-%s", resp)
		}
	}()
	var i int64
	for i = 1; i < 10; i++ {
		err := stream.Send(&protos.SendMessage{
			Name: strconv.FormatInt(i, 10),
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
	_ = stream.CloseSend()
	<-done
}
