package main

import (
	"context"
	"fmt"
	"log"
	"proto/quera/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	client, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	fsClient := pb.NewFSClient(client)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("api-key", "123456"))

	var (
		trailer metadata.MD
	)

	res, err := fsClient.Echo(ctx, &pb.EchoMessage{
		Msg: "Quera",
	}, grpc.Trailer(&trailer))

	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			log.Fatalf("error code : %d, msg : %s", status.Code(), status.Message())
		}
		log.Fatal(err)
	}

	fmt.Println(trailer.Get("x-user-id"))

	fmt.Println(res.GetEchoMsg())
}
