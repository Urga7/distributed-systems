package main

import (
	"net"
	"time"

	"homework04/chainStorage"
	"homework04/storage"

	"google.golang.org/grpc"
)

const lifetime time.Duration = 5 * time.Second

func Server(addr string, nextNodeAddr string, prevNodeAddr string) {
	todoStorage := storage.NewTodoStorage()
	chainNode := chainStorage.NewChainNode(todoStorage, nextNodeAddr, prevNodeAddr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		print("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	chainStorage.RegisterChainReplicationServer(grpcServer, chainNode)

	print("ChainNode server listening on %s", addr)
	go func() {
		time.Sleep(lifetime)
		print("Shutting down server...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(listener); err != nil {
		print("Failed to serve gRPC server: %v", err)
	}
}
