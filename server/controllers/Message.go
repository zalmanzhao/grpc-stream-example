package controllers

import (
	"context"
	"grpc.server/protos"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	protos.MessageServer
}

func (m *Message) Send(ctx context.Context, req *protos.SendMessage) (*protos.ReceiveMessage, error) {
	resp := &protos.ReceiveMessage{
		Name:    req.Name,
		Age:     18,
		Address: "beijing",
	}
	return resp, nil
}
func (m *Message) SendServerStream(req *protos.SendMessage, stream protos.Message_SendServerStreamServer) error {
	var i int64
	for i = 0; i < 10; i++ {
		resp := &protos.ReceiveMessage{
			Name:    req.Name + strconv.FormatInt(i, 10),
			Age:     18,
			Address: "beijing",
		}
		err := stream.Send(resp)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
func (m *Message) SendClientStream(stream protos.Message_SendClientStreamServer) error {
	var names []string
	for {
		resp := &protos.ReceiveMessage{
			Name:    "ClientStream Processing Completed" + strings.Join(names, ","),
			Age:     18,
			Address: "beijing",
		}
		res, err := stream.Recv()
		log.Println(names)
		if err == io.EOF {
			err := stream.SendAndClose(resp)
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		names = append(names, res.Name)
	}
}
func (m *Message) SendBidirectionalStream(stream protos.Message_SendClientServerStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		resp := &protos.ReceiveMessage{
			Name:    "ClientServerStream" + in.Name,
			Age:     18,
			Address: "beijing",
		}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
}
