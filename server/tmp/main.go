package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc.client/protos"
	"io"
	"log"
	"os"
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
	//sendServerStream(c)
	//sendClientStream(c)
	//sendBidirectionalStream(c)
	upload := protos.NewUploadClient(conn)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	sendClientUpload(upload)
}
func send(c protos.MessageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Send(ctx, &protos.SendMessage{Name: "SimpleAPI"})
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
		time.Sleep(1 * time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("failed to recv: %v", err)
	}
	log.Printf("Recv Ms: %s", resp)
}
func sendBidirectionalStream(c protos.MessageClient) {
	stream, err := c.SendBidirectionalStream(context.Background())
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
func sendClientUpload(c protos.UploadClient) {
	stream, err := c.Upload(context.Background())
	if err != nil {
		panic(err)
	}
	file := "main.go"
	fil, err := os.Open(file)
	if err != nil {
		return
	}

	// Maximum 1KB size per stream.
	buf := make([]byte, 1024)

	for {
		num, err := fil.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}

		if err := stream.Send(&protos.UploadRequest{Chunk: buf[:num], Name: file}); err != nil {
			log.Println(err)
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(resp)
	return
}
