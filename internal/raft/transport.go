package raft

import (
	"context"
	"time"

	pb "dist_db/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	rpcTimeout  = time.Second
	dialTimeout = 500 * time.Millisecond
)

func sendRequestVote(targetNode string, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), dialTimeout)
	defer dialCancel()

	connOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(dialCtx, targetNode, connOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	raftClient := pb.NewRaftClient(conn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer rpcCancel()

	return raftClient.RequestVote(rpcCtx, req)
}

func sendAppendEntries(targetNode string, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), dialTimeout)
	defer dialCancel()

	connOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	conn, err := grpc.DialContext(dialCtx, targetNode, connOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	raftClient := pb.NewRaftClient(conn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer rpcCancel()

	return raftClient.AppendEntries(rpcCtx, req)
}
