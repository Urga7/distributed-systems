package chainStorage

import (
	"homework04/storage"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type ChainNode struct {
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

func (s *ChainNode) Put(ctx context.Context, putRequest *PutRequest) *PutResponse {
	s.todoStore.Put(putRequest.Key, putRequest.Value)

	// If this node is not the tail, forward the request to the next node
	if s.nextNodeAddr != "" {
		conn, err := grpc.Dial(s.nextNodeAddr, grpc.WithInsecure())
		if err != nil {
			return nil
		}
		defer conn.Close()

		client := NewChainReplicationClient(conn)
		_, err = client.Put(ctx, putRequest)
		if err != nil {
			return nil
		}
	}

	return &PutResponse{Status: "OK"}
}

func (s *ChainNode) Get(ctx context.Context, getRequest *GetRequest) *GetResponse {
	value, found := s.todoStore.Get(getRequest.Key)
	if !found {
		return &GetResponse{Value: "", Found: false}
	}

	return &GetResponse{Value: value, Found: true}
}

func (s *ChainNode) Commit(ctx context.Context, todo *Todo) *PutResponse {
	// Commit the key locally
	err := s.todoStore.Commit(todo.Key)
	if err != nil {
		return nil
	}

	// If this node is not the head, propagate the commit to the previous node
	if s.prevNodeAddr != "" {
		conn, err := grpc.Dial(s.prevNodeAddr, grpc.WithInsecure())
		if err != nil {
			return nil
		}
		defer conn.Close()

		client := NewChainReplicationClient(conn)
		_, err = client.Commit(ctx, todo)
		if err != nil {
			return nil
		}
	}

	return &PutResponse{Status: "Committed"}
}
