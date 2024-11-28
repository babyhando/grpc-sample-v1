package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"proto/quera/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var filePath = flag.String("file", "", "file to send")

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatal("no file provided to upload")
	}

	client, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	fsClient := pb.NewFSClient(client)

	upload(fsClient)

	// download(fsClient, *filePath)
}

func upload(fsClient pb.FSClient) {
	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Base(*filePath)

	stream, err := fsClient.Upload(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	start, end := 0, 500
	chunkSize := 500
	for {
		if end >= len(data) {
			end = len(data)
		}

		done := end == len(data)

		dataToSend := data[start:end]

		stream.Send(&pb.Chunk{
			Data:     dataToSend,
			Done:     done,
			FileName: filename,
		})

		if done {
			break
		}

		start, end = end, end+chunkSize
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UPLOAD STATUS : ", resp.Status)
}

func download(fsClient pb.FSClient, filePath string) {
	ctx := context.Background()

	stream, err := fsClient.Download(ctx, &pb.DownloadRequest{
		FileName: filePath,
	})

	if err != nil {
		log.Fatal(err)
	}

	var data []byte

	for {
		chunk, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, chunk.GetData()...)

		if chunk.Done {
			break
		}
	}

	if err := os.WriteFile("downloaded-"+filePath, data, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
