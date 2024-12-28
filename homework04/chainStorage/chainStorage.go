package chainStorage

import (
	"fmt"
	"homework04/storage"
	"time"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type ChainNode struct {
	UnimplementedChainReplicationServer
	todoStore    *storage.TodoStorage
	nextNodeAddr string
	prevNodeAddr string
}

func NewChainNode(todoStorage *storage.TodoStorage, nextNodeAddr, prevNodeAddr string) *ChainNode {
	return &ChainNode{
		todoStore:    todoStorage,
		nextNodeAddr: nextNodeAddr,
		prevNodeAddr: prevNodeAddr,
	}
}

func currentTimestamp() string {
	return time.Now().Format("15:04:05.000")
}

func print(format string, args ...interface{}) {
	fmt.Printf("[%s] "+format+"\n", append([]interface{}{currentTimestamp()}, args...)...)
}

func (s *ChainNode) Put(ctx context.Context, putRequest *PutRequest) (*PutResponse, error) {
	print("Received Put request.")
	s.todoStore.Put(putRequest.Key, putRequest.Value)
	print("Entry saved locally: [%s] = %s", putRequest.Key, putRequest.Value)

	if s.nextNodeAddr == "" {
		err := s.todoStore.Commit(putRequest.Key)
		if err != nil {
			return nil, err
		}
		print("Entry committed at chain tail: [%s] = %s", putRequest.Key, putRequest.Value)
	} else {
		conn, err := grpc.Dial(s.nextNodeAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := NewChainReplicationClient(conn)

		print("Forwarding Put request to %s", s.nextNodeAddr)
		_, err = client.Put(ctx, putRequest)
		if err != nil {
			return nil, err
		}

		err = s.todoStore.Commit(putRequest.Key)
		if err != nil {
			return nil, err
		}
		print("Received Put confirmation from %s, entry committed: [%s] = %s", s.nextNodeAddr, putRequest.Key, putRequest.Value)
	}

	return &PutResponse{Status: "OK"}, nil
}

func (s *ChainNode) Get(ctx context.Context, getRequest *GetRequest) (*GetResponse, error) {
	print("Received Get request for key: %s", getRequest.Key)
	value, found, commited := s.todoStore.Get(getRequest.Key)
	if !found {
		print("Entry not found for key: %s", getRequest.Key)
		return &GetResponse{Value: "", Found: false}, nil
	}

	if !commited {
		print("Found uncommited entry for key: %s, retrying in 0.1 seconds", getRequest.Key)
		time.Sleep(100 * time.Millisecond)
		return s.Get(ctx, getRequest)
	}

	print("Found commited entry for key: %s", getRequest.Key)
	return &GetResponse{Value: value, Found: true}, nil
}
