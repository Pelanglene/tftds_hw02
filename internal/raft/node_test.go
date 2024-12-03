package raft_test

import (
	pb "dist_db/api/proto"
	"dist_db/internal/raft"
	server "dist_db/internal/server"
	"fmt"
	"log"
	"log/slog"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
)

const (
	BeginTestRaftPort = 5050
	BeginTestHttpPort = 8080
)

type TestRaftServer struct {
	raftServer *raft.RaftServer
	grpcServer *grpc.Server // for stops
	raftPort   string

	httpServer *server.HTTPServer
	httpPort   string
}

func NewTestServer(id int64, peers []string, raftPort, httpPort int64) *TestRaftServer {
	raftServer := raft.NewRaftServer(id, peers)
	httpServer := server.NewHTTPServer(raftServer)

	return &TestRaftServer{
		raftServer: raftServer,
		httpServer: httpServer,

		raftPort: fmt.Sprintf(":%d", raftPort),
		httpPort: fmt.Sprintf(":%d", httpPort),
	}
}

func StartTestServer(t *TestRaftServer) {
	grpcServer := grpc.NewServer()
	pb.RegisterRaftServer(grpcServer, t.raftServer)
	t.grpcServer = grpcServer

	t.raftServer.StartTimeouts()

	lis, err := net.Listen("tcp", t.raftPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		slog.Info("Raft-server starting", "port", t.raftPort)
		if err := t.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}

func StopTestServer(t *TestRaftServer) {
	t.grpcServer.GracefulStop()
	t.raftServer.ResetTimeouts()
}

func Filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func NewTestCluster(count int) []*TestRaftServer {
	result := make([]*TestRaftServer, 0, count)

	peers := []string{}
	for id := 0; id < count; id++ {
		peers = append(peers, fmt.Sprintf("localhost:%d", BeginTestRaftPort+id))
	}

	for id := int64(0); id < int64(count); id++ {
		newServer := NewTestServer(
			id,
			Filter(peers, func(s string) bool {
				return s != fmt.Sprintf("localhost:%d", BeginTestRaftPort+id)
			}),
			BeginTestRaftPort+id,
			BeginTestHttpPort+id,
		)

		StartTestServer(newServer)
		result = append(result, newServer)
	}

	return result
}

func TestClusterOperations(t *testing.T) {
	cluster := NewTestCluster(3)
	defer func() {
		for _, server := range cluster {
			StopTestServer(server)
		}
	}()

	time.Sleep(10 * time.Second)

	var leader *TestRaftServer
	for _, server := range cluster {
		if server.raftServer.State == raft.LEADER {
			leader = server
			break
		}
	}

	if leader == nil {
		t.Fatal("No leader elected")
	}

	t.Run("Write Operation", func(t *testing.T) {
		value := "test_value"
		success, err := leader.raftServer.ReplicateLogEntry("CREATE", "test_key", &value, nil)
		if err != nil {
			t.Errorf("Write operation failed: %v", err)
		}
		if !success {
			t.Error("Write operation should succeed")
		}
	})

	t.Run("Replication Check", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		for _, server := range cluster {
			if len(server.raftServer.Log) < 2 {
				t.Errorf("Log not replicated to node %d", server.raftServer.ID)
			}
		}
	})
}

func TestLeaderElection(t *testing.T) {
	cluster := NewTestCluster(5)
	defer func() {
		for _, server := range cluster {
			StopTestServer(server)
		}
	}()

	time.Sleep(10 * time.Second)

	leaders := 0
	for _, server := range cluster {
		if server.raftServer.State == raft.LEADER {
			leaders++
		}
	}

	if leaders != 1 {
		t.Errorf("Expected exactly one leader, got %d", leaders)
	}
}

func TestLeaderFailure(t *testing.T) {
	cluster := NewTestCluster(5)
	defer func() {
		for _, server := range cluster {
			StopTestServer(server)
		}
	}()

	time.Sleep(10 * time.Second)

	var leader *TestRaftServer
	var followers []*TestRaftServer
	for _, server := range cluster {
		if server.raftServer.State == raft.LEADER {
			leader = server
		} else {
			followers = append(followers, server)
		}
	}

	if leader == nil {
		t.Fatal("No leader elected")
	}

	value := "test_value"
	success, err := leader.raftServer.ReplicateLogEntry("CREATE", "test_key", &value, nil)
	if err != nil {
		t.Errorf("Write operation failed: %v", err)
	}
	if !success {
		t.Error("Write operation should succeed")
	}

	time.Sleep(2 * time.Second)

	initialLogSize := len(leader.raftServer.Log)
	for _, server := range followers {
		if len(server.raftServer.Log) != initialLogSize {
			t.Errorf("Log not replicated before leader failure. Expected size %d, got %d",
				initialLogSize, len(server.raftServer.Log))
		}
	}

	StopTestServer(leader)
	time.Sleep(15 * time.Second)

	newLeaderCount := 0
	var newLeader *TestRaftServer
	for _, server := range followers {
		if server.raftServer.State == raft.LEADER {
			newLeaderCount++
			newLeader = server
		}
	}

	if newLeaderCount != 1 {
		t.Errorf("Expected exactly one new leader, got %d", newLeaderCount)
	}

	if newLeader == nil {
		t.Fatal("No new leader elected after failure")
	}

	// Проверяем, что новый лидер может принимать записи
	newValue := "test_value_after_failure"
	success, err = newLeader.raftServer.ReplicateLogEntry("CREATE", "test_key_2", &newValue, nil)
	if err != nil {
		t.Errorf("Write operation failed with new leader: %v", err)
	}
	if !success {
		t.Error("Write operation with new leader should succeed")
	}

	time.Sleep(2 * time.Second)

	expectedLogSize := initialLogSize + 1
	for _, server := range followers {
		if len(server.raftServer.Log) != expectedLogSize {
			t.Errorf("Log not replicated after leader failure. Expected size %d, got %d",
				expectedLogSize, len(server.raftServer.Log))
		}
	}
}

func TestNetworkPartition(t *testing.T) {
	cluster := NewTestCluster(7)
	defer func() {
		for _, server := range cluster {
			StopTestServer(server)
		}
	}()

	time.Sleep(10 * time.Second)

	var leader *TestRaftServer
	var majorityGroup, minorityGroup []*TestRaftServer

	for _, server := range cluster {
		if server.raftServer.State == raft.LEADER {
			leader = server
			majorityGroup = append(majorityGroup, server)
		} else {
			if len(majorityGroup) < 4 {
				majorityGroup = append(majorityGroup, server)
			} else {
				minorityGroup = append(minorityGroup, server)
			}
		}
	}

	if leader == nil {
		t.Fatal("No leader elected")
	}

	t.Run("Initial State", func(t *testing.T) {
		value := "test_value"
		success, err := leader.raftServer.ReplicateLogEntry("CREATE", "test_key", &value, nil)
		if err != nil {
			t.Errorf("Initial write failed: %v", err)
		}
		if !success {
			t.Error("Initial write should succeed")
		}
		time.Sleep(2 * time.Second)
	})

	// Симулируем разделение сети
	t.Run("Network Partition", func(t *testing.T) {
		for _, majorityServer := range majorityGroup {
			for _, minorityServer := range minorityGroup {
				majorityServer.raftServer.RemovePeer(minorityServer.raftPort[1:])
				minorityServer.raftServer.RemovePeer(majorityServer.raftPort[1:])
			}
		}

		// Проверяем, что операции записи все еще работают в группе большинства
		value := "value_after_partition"
		success, err := leader.raftServer.ReplicateLogEntry("CREATE", "key_after_partition", &value, nil)
		if err != nil {
			t.Errorf("Write during partition failed: %v", err)
		}
		if !success {
			t.Error("Write during partition should succeed in majority group")
		}

		time.Sleep(2 * time.Second)

		// Проверяем, что запись реплицировалась только в группу большинства
		majorityLogSize := len(leader.raftServer.Log)

		for _, server := range majorityGroup {
			if len(server.raftServer.Log) != majorityLogSize {
				t.Errorf("Log not replicated in majority group. Server %d has size %d, expected %d",
					server.raftServer.ID, len(server.raftServer.Log), majorityLogSize)
			}
		}
	})

	// Восстанавливаем сетевое соединение
	t.Run("Network Healing", func(t *testing.T) {
		for _, majorityServer := range majorityGroup {
			for _, minorityServer := range minorityGroup {
				majorityServer.raftServer.AddPeer(fmt.Sprintf("localhost%s", minorityServer.raftPort))
				minorityServer.raftServer.AddPeer(fmt.Sprintf("localhost%s", majorityServer.raftPort))
			}
		}

		time.Sleep(5 * time.Second)

		value := "value_after_healing"
		success, err := leader.raftServer.ReplicateLogEntry("CREATE", "key_after_healing", &value, nil)
		if err != nil {
			t.Errorf("Write after healing failed: %v", err)
		}
		if !success {
			t.Error("Write after healing should succeed")
		}

		time.Sleep(2 * time.Second)

		// Проверяем, что все узлы имеют одинаковый размер лога
		expectedLogSize := len(leader.raftServer.Log)
		for _, server := range cluster {
			if len(server.raftServer.Log) != expectedLogSize {
				t.Errorf("After healing: Server %d has incorrect log size. Got %d, expected %d",
					server.raftServer.ID, len(server.raftServer.Log), expectedLogSize)
			}
		}

		// Проверяем, что последняя запись есть на всех узлах
		lastEntry := leader.raftServer.Log[len(leader.raftServer.Log)-1]
		for _, server := range cluster {
			serverLastEntry := server.raftServer.Log[len(server.raftServer.Log)-1]
			if serverLastEntry.Key != lastEntry.Key ||
				*serverLastEntry.Value != *lastEntry.Value {
				t.Errorf("Server %d has incorrect last entry. Expected key=%v value=%v, got key=%v value=%v",
					server.raftServer.ID,
					lastEntry.Key,
					*lastEntry.Value,
					serverLastEntry.Key,
					*serverLastEntry.Value)
			}
		}
	})
}
