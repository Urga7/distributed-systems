package main

import (
	"context"
	"fmt"
	"time"

	"homework04/chainStorage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func currentTimestamp() string {
	return time.Now().Format("15:04:05.000")
}

func print(format string, args ...interface{}) {
	fmt.Printf("[%s] "+format+"\n", append([]interface{}{currentTimestamp()}, args...)...)
}

func Client(putUrl string, getUrl string) {
	putConn, err := grpc.Dial(putUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		print("Failed to connect to Put server: %v", err)
		panic(err)
	}
	defer putConn.Close()

	print("Client connected to Put server at %s", putUrl)
	putClient := chainStorage.NewChainReplicationClient(putConn)

	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("task%d", i)
		value := fmt.Sprintf("Task %d description", i)
		putReq := &chainStorage.PutRequest{Key: key, Value: value}

		_, err := putClient.Put(context.Background(), putReq)
		if err != nil {
			print("Put request failed: %v", err)
			panic(err)
		}
		print("Put: Key = %s, Value = %s", key, value)
	}

	getConn, err := grpc.Dial(getUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		print("Failed to connect to Get server: %v", err)
		panic(err)
	}
	defer getConn.Close()

	print("Client connected to Get server at %s", getUrl)
	getClient := chainStorage.NewChainReplicationClient(getConn)

	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("task%d", i)
		getReq := &chainStorage.GetRequest{Key: key}

		resp, err := getClient.Get(context.Background(), getReq)
		if err != nil {
			print("Get request failed: %v", err)
			panic(err)
		}

		if resp.Found {
			print("Get: Key = %s, Value = %s", key, resp.Value)
		} else {
			print("Get: Key = %s not found", key)
		}
	}
}
