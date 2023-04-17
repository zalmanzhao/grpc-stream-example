package controllers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc.server/protos"
	"grpc.server/storage"
	"io"
)

type Upload struct {
	protos.UploadServer
	storage storage.Manager
}

func NewUploadServer(storage storage.Manager) Upload {
	return Upload{
		storage: storage,
	}
}

func (s Upload) Upload(stream protos.Upload_UploadServer) error {
	var name string
	var file *storage.File
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := s.storage.Store(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			return stream.SendAndClose(&protos.UploadResponse{Name: name})
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		name = req.Name
		if file == nil {
			file = storage.NewFile(name)
		}
		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
