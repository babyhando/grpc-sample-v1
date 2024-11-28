package api

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"proto/quera/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fsGRPCApi struct {
	pb.UnimplementedFSServer
	fileDir string
}

func NewFsGRPCApi(fileDir string) pb.FSServer {
	return &fsGRPCApi{fileDir: fileDir}
}

func (s *fsGRPCApi) Upload(req grpc.ClientStreamingServer[pb.Chunk, pb.UploadResponse]) error {
	var (
		filename string
		content  []byte
	)

	// handle stream
	for {
		chunk, err := req.Recv()
		if err != nil {
			log.Println("error on recv chunk ", err.Error())
			return req.SendAndClose(&pb.UploadResponse{
				Status: pb.UploadStatus_FAILED,
			})
		}

		if len(filename) == 0 {
			filename = chunk.GetFileName()
		}

		content = append(content, chunk.Data...)

		if chunk.Done {
			break
		}
	}

	// save to disk
	if err := os.WriteFile(filepath.Join(s.fileDir, filename), content, os.ModePerm); err != nil {
		log.Println("error on save file to disk : ", err.Error())
		return req.SendAndClose(&pb.UploadResponse{
			Status: pb.UploadStatus_FAILED,
		})
	}

	return req.SendAndClose(&pb.UploadResponse{
		Status: pb.UploadStatus_SUCCESS,
	})
}

func (s *fsGRPCApi) Download(req *pb.DownloadRequest, stream grpc.ServerStreamingServer[pb.Chunk]) error {
	data, err := os.ReadFile(filepath.Join(s.fileDir, req.GetFileName()))
	if err != nil {
		return err
	}

	chunkSize := 512
	start, end := 0, chunkSize
	for {
		if end >= len(data) {
			end = len(data)
		}

		done := end == len(data)

		dataToSend := data[start:end]

		stream.Send(&pb.Chunk{
			Data: dataToSend,
			Done: done,
		})

		if done {
			return nil
		}

		start, end = end, end+chunkSize
	}
}

func (s *fsGRPCApi) Echo(ctx context.Context, req *pb.EchoMessage) (*pb.EchoResponse, error) {
	grpc.SetTrailer(ctx, metadata.Pairs("x-user-id", "2324"))
	resp := &pb.EchoResponse{
		EchoMsg: "echoed from server : " + req.GetMsg(),
	}

	return resp, nil
}

func (s *fsGRPCApi) Echo2(ctx context.Context, req *pb.EchoMessage) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		EchoMsg: "echoed 2 from server : " + req.GetMsg(),
	}, nil
}
