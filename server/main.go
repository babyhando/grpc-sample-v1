package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"proto/quera/api"
	"proto/quera/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logHandler), grpc.StreamInterceptor(logStreamHandler))

	pb.RegisterOrderServiceServer(grpcServer, api.NewOrderGRPCServer())
	pb.RegisterFSServer(grpcServer, api.NewFsGRPCApi("./files"))

	log.Println("listening on :8080 ...")
	grpcServer.Serve(l)
}

var logHandler grpc.UnaryServerInterceptor = func(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (resp any, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "header should be exists")
	}

	if v := md.Get("api-key"); len(v) == 0 || v[0] != "123456" {
		return nil, status.Errorf(codes.Unauthenticated, "wrong api key")
	}

	fmt.Printf("Unary Req Type : %T - Method : %s\n", req, info.FullMethod)

	res, err := next(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Unary Res Type : %T", res)

	return res, nil
}

var logStreamHandler grpc.StreamServerInterceptor = func(srv any, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, next grpc.StreamHandler) error {

	fmt.Printf("Server Handler Type : %T - Method : %s\n", srv, info.FullMethod)

	if err := next(srv, ss); err != nil {
		log.Println("error on server stream", err.Error())
		return err
	}

	// server stream had finished

	return nil
}
