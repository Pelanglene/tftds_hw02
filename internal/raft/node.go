package raft

import (
	"context"
	pb "dist_db/api/proto"
	"dist_db/internal/db"
	"fmt"
	"sync/atomic"

	"log"
	"net"
	"sync"
	"time"

	"log/slog"

	"google.golang.org/grpc"
)

const (
	LEADER = iota
	FOLLOWER
	CANDIDATE
)

type LogEntry struct {
	Term     int64
	Command  string
	Key      string
	Value    *string
	OldValue *string
}

type TimeoutManager struct {
	electionDuration  time.Duration
	heartbeatDuration time.Duration

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer

	mu sync.Mutex
}

func NewTimeoutManager(nodeID int64) *TimeoutManager {
	return &TimeoutManager{
		electionDuration:  time.Second * time.Duration(8+nodeID*5),
		heartbeatDuration: time.Second * 5,
	}
}

func (tm *TimeoutManager) scheduleTimeout(timer **time.Timer, duration time.Duration, callback func()) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if *timer != nil {
		(*timer).Stop()
	}
	*timer = time.AfterFunc(duration, callback)
}

func (tm *TimeoutManager) StartElectionTimer(callback func()) {
	tm.scheduleTimeout(&tm.electionTimer, tm.electionDuration, callback)
}

func (tm *TimeoutManager) StartHeartbeatTimer(callback func()) {
	tm.scheduleTimeout(&tm.heartbeatTimer, tm.heartbeatDuration, callback)
}

func (tm *TimeoutManager) StopAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.electionTimer != nil {
		tm.electionTimer.Stop()
	}
	if tm.heartbeatTimer != nil {
		tm.heartbeatTimer.Stop()
	}
}

type RaftServer struct {
	pb.UnimplementedRaftServer

	ID           int64
	currentTerm  int64
	lastVotedFor int64
	Log          []LogEntry

	State       int
	leaderID    int64
	commitIndex int64
	nextIndex   map[string]int64

	timeouts *TimeoutManager
	peers    []string
	mu       sync.Mutex
}

func NewRaftServer(id int64, peers []string) *RaftServer {
	server := &RaftServer{
		ID:           id,
		currentTerm:  0,
		lastVotedFor: -1,
		Log: []LogEntry{
			{
				Term:    0,
				Command: "init",
			},
		},
		State:       FOLLOWER,
		leaderID:    -1,
		commitIndex: 0,
		nextIndex:   make(map[string]int64),
		timeouts:    NewTimeoutManager(id),
		peers:       peers,
	}

	server.timeouts.StartElectionTimer(server.beginElection)

	return server
}

func (s *RaftServer) RequestVote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	slog.Debug("Received vote request",
		"node", s.ID,
		"candidate", req.CandidateID,
		"candidate_term", req.Term,
		"current_term", s.currentTerm)

	if !s.isTermValid(req.Term) {
		return &pb.VoteResponse{
			Term:        s.currentTerm,
			VoteGranted: false,
		}, nil
	}

	canVote := s.canVoteFor(req.CandidateID, req.Term)
	if canVote {
		s.updateStateForVote(req.Term, req.CandidateID)
		slog.Info("Vote granted",
			"node", s.ID,
			"candidate", req.CandidateID,
			"term", req.Term)
	}

	return &pb.VoteResponse{
		Term:        s.currentTerm,
		VoteGranted: canVote,
	}, nil
}

func (s *RaftServer) isTermValid(requestTerm int64) bool {
	return requestTerm >= s.currentTerm
}

func (s *RaftServer) canVoteFor(candidateID int64, requestTerm int64) bool {
	return s.lastVotedFor == -1 ||
		requestTerm > s.currentTerm ||
		s.lastVotedFor == candidateID
}

func (s *RaftServer) updateStateForVote(newTerm int64, candidateID int64) {
	s.currentTerm = newTerm
	s.lastVotedFor = candidateID
}

func (s *RaftServer) appendEntries(req *pb.AppendEntriesRequest) {
	for idx, newEntry := range req.Entries {
		targetIndex := req.PrevLogIndex + int64(idx) + 1

		if s.hasConflictingEntry(targetIndex, req.Term) {
			s.handleLogConflict(targetIndex)
			continue
		}

		s.appendNewEntry(newEntry)
	}
}

func (s *RaftServer) hasConflictingEntry(index int64, term int64) bool {
	return index < int64(len(s.Log)) && s.Log[index].Term != term
}

func (s *RaftServer) handleLogConflict(index int64) {
	s.Log = s.Log[:index]
}

func (s *RaftServer) appendNewEntry(entry *pb.LogEntry) {
	newLogEntry := LogEntry{
		Term:     entry.Term,
		Command:  entry.Command,
		Key:      entry.Key,
		Value:    entry.Value,
		OldValue: entry.OldValue,
	}
	s.Log = append(s.Log, newLogEntry)
}

func (s *RaftServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	slog.Info("AppendEntries received", "node", s.ID, "leader", req.LeaderID)

	if req.Term < s.currentTerm {
		return &pb.AppendEntriesResponse{Term: s.currentTerm, Success: false}, nil
	}

	s.currentTerm = req.Term
	s.lastVotedFor = -1
	s.State = FOLLOWER
	s.leaderID = req.LeaderID
	s.timeouts.StartElectionTimer(s.beginElection)

	if req.PrevLogIndex >= 0 {
		if req.PrevLogIndex >= int64(len(s.Log)) || s.Log[req.PrevLogIndex].Term != req.PrevLogTerm {
			return &pb.AppendEntriesResponse{Term: s.currentTerm, Success: false}, nil
		}
	}

	s.appendEntries(req)

	if req.LeaderCommit > s.commitIndex {
		s.applyEntries(req.LeaderCommit)
	}

	return &pb.AppendEntriesResponse{Term: s.currentTerm, Success: true}, nil
}

func (s *RaftServer) applyEntries(leaderCommit int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.commitIndex >= leaderCommit {
		slog.Debug("Nothing to apply: current index is greater or equal to leader's index",
			"node", s.ID,
			"current_index", s.commitIndex,
			"leader_index", leaderCommit)
		return
	}

	slog.Info("Start applying entries",
		"node", s.ID,
		"current_index", s.commitIndex,
		"leader_index", leaderCommit)

	for currentIndex := s.commitIndex + 1; currentIndex <= leaderCommit; currentIndex++ {
		if currentIndex >= int64(len(s.Log)) {
			slog.Error("Applied index exceeds log size",
				"node", s.ID,
				"index", currentIndex,
				"log_size", len(s.Log))
			break
		}

		entry := s.Log[currentIndex]

		if entry.Command == "init" {
			slog.Debug("Skip init entry",
				"node", s.ID,
				"index", currentIndex)
			continue
		}

		db.ProcessWrite(entry.Command, entry.Key, entry.Value, entry.OldValue)
		// if err != nil {
		// 	slog.Error("Error applying entry",
		// 		"node", s.ID,
		// 		"index", currentIndex,
		// 		"entry", entry,
		// 		"error", err)
		// 	continue
		// }

		// if !success {
		// 	slog.Warn("Failed to apply entry",
		// 		"node", s.ID,
		// 		"index", currentIndex,
		// 		"entry", entry)
		// 	continue
		// }

		slog.Info("Successfully applied entry",
			"node", s.ID,
			"index", currentIndex,
			"entry", entry)
	}

	s.commitIndex = leaderCommit
	slog.Info("Updated applied entries index",
		"node", s.ID,
		"new_index", s.commitIndex)
}

func (s *RaftServer) StartRaftServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRaftServer(grpcServer, s)

	slog.Info("Raft server starts to listen", "node_id", s.ID, "port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *RaftServer) GetLeaderID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.leaderID
}

func (s *RaftServer) ResetTimeouts() {
	s.timeouts.StopAll()
}

func (s *RaftServer) StartTimeouts() {
	s.timeouts.StartElectionTimer(s.beginElection)
	if s.State == LEADER {
		s.timeouts.StartHeartbeatTimer(s.sendHeartbeats)
	}
}

func (s *RaftServer) LogLength() int {
	return len(s.Log)
}

func (s *RaftServer) beginElection() {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("Begin election...",
		"node", s.ID,
		"current_term", s.currentTerm)

	s.promoteToCandidate()

	voteRequest := s.prepareVoteRequest()
	receivedVotes := &atomicVoteCounter{count: 1} // Голосуем за себя
	requiredVotes := len(s.peers)/2 + 1

	s.timeouts.StartElectionTimer(s.beginElection)

	s.requestVotesFromPeers(voteRequest, receivedVotes, requiredVotes)
}

type atomicVoteCounter struct {
	count int32
}

func (s *RaftServer) promoteToCandidate() {
	s.currentTerm++
	s.State = CANDIDATE
	s.lastVotedFor = s.ID
	s.leaderID = -1
}

func (s *RaftServer) prepareVoteRequest() *pb.VoteRequest {
	lastLogIndex := int64(len(s.Log) - 1)
	var lastLogTerm int64
	if lastLogIndex >= 0 {
		lastLogTerm = s.Log[lastLogIndex].Term
	}

	return &pb.VoteRequest{
		Term:         s.currentTerm,
		CandidateID:  s.ID,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}
}

func (s *RaftServer) requestVotesFromPeers(request *pb.VoteRequest, votes *atomicVoteCounter, required int) {
	for _, peer := range s.peers {
		go func(peer string) {
			response, err := sendRequestVote(peer, request)

			if err != nil {
				slog.Error("Error requesting vote",
					"node", s.ID,
					"peer", peer,
					"error", err)
				return
			}

			if response.VoteGranted {
				s.handleVoteGranted(votes, required)
			} else {
				slog.Debug("Vote not granted",
					"node", s.ID,
					"peer", peer,
					"response_term", response.Term)

				if response.Term > request.Term {
					s.mu.Lock()
					defer s.mu.Unlock()
					s.currentTerm = response.Term
					s.State = FOLLOWER
					s.lastVotedFor = -1
				}
			}
		}(peer)
	}
}

func (s *RaftServer) handleVoteGranted(votes *atomicVoteCounter, required int) {
	newCount := atomic.AddInt32(&votes.count, 1)

	slog.Info("Vote granted!",
		"node", s.ID,
		"total_votes", newCount,
		"required", required)

	if int(newCount) >= required {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.State == CANDIDATE {
			s.becomeLeader()
		}
	}
}

func (s *RaftServer) becomeLeader() {
	slog.Info("Become leader", "node", s.ID)

	s.State = LEADER
	s.leaderID = s.ID

	s.timeouts.StopAll()
	s.timeouts.StartHeartbeatTimer(s.sendHeartbeats)
}

func (s *RaftServer) sendHeartbeats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State != LEADER {
		return
	}

	for _, peer := range s.peers {
		go s.synchronizePeerLog(peer)
	}

	s.timeouts.StartHeartbeatTimer(s.sendHeartbeats)
}

func (s *RaftServer) synchronizePeerLog(peer string) {
	for {
		entries, nextIndex, prevLogTerm := s.prepareEntriesForPeer(peer)
		if nextIndex == 0 {
			return
		}

		req := s.createHeartbeatRequest(nextIndex, prevLogTerm, entries)
		success := s.sendHeartbeatRequest(peer, req, nextIndex, len(entries))

		if success {
			break
		}
	}
}

func (s *RaftServer) prepareEntriesForPeer(peer string) ([]*pb.LogEntry, int64, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextIndex, exists := s.nextIndex[peer]
	if !exists {
		s.nextIndex[peer] = 0
		return nil, 0, 0
	}

	entries := s.Log[nextIndex:]
	entriesProto := make([]*pb.LogEntry, len(entries))

	for i, entry := range entries {
		entriesProto[i] = &pb.LogEntry{
			Term:     entry.Term,
			Command:  entry.Command,
			Key:      entry.Key,
			Value:    entry.Value,
			OldValue: entry.OldValue,
		}
	}

	var prevLogTerm int64
	if nextIndex > 0 {
		prevLogTerm = s.Log[nextIndex-1].Term
	}

	return entriesProto, nextIndex, prevLogTerm
}

func (s *RaftServer) createHeartbeatRequest(nextIndex int64, prevLogTerm int64, entries []*pb.LogEntry) *pb.AppendEntriesRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderID:     s.ID,
		LeaderCommit: s.commitIndex,
		PrevLogIndex: nextIndex - 1,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
	}
}

func (s *RaftServer) sendHeartbeatRequest(peer string, req *pb.AppendEntriesRequest, nextIndex int64, entriesCount int) bool {
	resp, err := sendAppendEntries(peer, req)
	if err != nil {
		slog.Error("Error sending heartbeat",
			"leader", s.ID,
			"peer", peer,
			"error", err)
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if resp.Success {
		s.nextIndex[peer] = nextIndex + int64(entriesCount)
		return true
	}

	s.nextIndex[peer] = nextIndex - 1
	slog.Info("Node is not synchronized! Decrease nextIndex and repeat",
		"leader", s.ID,
		"node", peer,
		"new_index", nextIndex-1)
	return false
}

func (s *RaftServer) ReplicateLogEntry(command, key string, value, oldValue *string) (bool, error) {
	s.mu.Lock()

	if s.State != LEADER {
		s.mu.Unlock()
		return false, fmt.Errorf("cannot handle write on replica. LeaderID: %v", s.leaderID)
	}

	logEntry := s.prepareLogEntry(command, key, value, oldValue)
	replicationInfo := s.appendAndPrepareReplication(logEntry)
	s.mu.Unlock()

	slog.Info("Leader append new entry",
		"node", s.ID,
		"entry", logEntry)

	success := s.replicateToFollowers(replicationInfo)
	if !success {
		return false, fmt.Errorf("не удалось реплицировать запись %+v", logEntry)
	}

	return s.applyReplicatedEntry(logEntry, replicationInfo.currentIndex)
}

type replicationInfo struct {
	prevLogIndex  int64
	prevLogTerm   int64
	currentIndex  int64
	peers         []string
	successCount  *int32
	requiredVotes int32
}

func (s *RaftServer) prepareLogEntry(command, key string, value, oldValue *string) LogEntry {
	return LogEntry{
		Term:     s.currentTerm,
		Command:  command,
		Key:      key,
		Value:    value,
		OldValue: oldValue,
	}
}

func (s *RaftServer) appendAndPrepareReplication(entry LogEntry) *replicationInfo {
	prevLogIndex := int64(len(s.Log) - 1)
	var prevLogTerm int64
	if prevLogIndex >= 0 {
		prevLogTerm = s.Log[prevLogIndex].Term
	}

	s.Log = append(s.Log, entry)
	currentIndex := int64(len(s.Log) - 1)

	peers := make([]string, len(s.peers))
	copy(peers, s.peers)

	return &replicationInfo{
		prevLogIndex:  prevLogIndex,
		prevLogTerm:   prevLogTerm,
		currentIndex:  currentIndex,
		peers:         peers,
		successCount:  new(int32),
		requiredVotes: int32(len(peers)/2 + 1),
	}
}

func (s *RaftServer) replicateToFollowers(info *replicationInfo) bool {
	atomic.StoreInt32(info.successCount, 1) // Считаем голос лидера

	var wg sync.WaitGroup
	for _, peer := range info.peers {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()
			s.replicateToPeer(peer, info)
		}(peer)
	}
	wg.Wait()

	return atomic.LoadInt32(info.successCount) >= info.requiredVotes
}

func (s *RaftServer) replicateToPeer(peer string, info *replicationInfo) {
	req := s.createReplicationRequest(info)
	resp, err := sendAppendEntries(peer, req)

	if err == nil && resp.Success {
		atomic.AddInt32(info.successCount, 1)
		s.mu.Lock()
		s.nextIndex[peer] = info.currentIndex + 1
		s.mu.Unlock()
	}
}

func (s *RaftServer) createReplicationRequest(info *replicationInfo) *pb.AppendEntriesRequest {
	lastEntry := s.Log[info.currentIndex]
	return &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderID:     s.ID,
		PrevLogIndex: info.prevLogIndex,
		PrevLogTerm:  info.prevLogTerm,
		Entries: []*pb.LogEntry{{
			Term:     lastEntry.Term,
			Command:  lastEntry.Command,
			Key:      lastEntry.Key,
			Value:    lastEntry.Value,
			OldValue: lastEntry.OldValue,
		}},
		LeaderCommit: s.commitIndex,
	}
}

func (s *RaftServer) applyReplicatedEntry(entry LogEntry, currentIndex int64) (bool, error) {
	s.mu.Lock()
	s.commitIndex = currentIndex
	s.mu.Unlock()

	slog.Info("Successfully replicated entry, applying...",
		"leader", s.ID,
		"entry", entry)

	return db.ProcessWrite(entry.Command, entry.Key, entry.Value, entry.OldValue)
}

type Replicator struct {
	raftServer *RaftServer
}

func NewReplicator(raftServer *RaftServer) *Replicator {
	return &Replicator{
		raftServer: raftServer,
	}
}

func (r *Replicator) ApplyAndReplicate(command string, key string, value, oldValue *string) (bool, error) {
	return r.raftServer.ReplicateLogEntry(command, key, value, oldValue)
}

func (s *RaftServer) RemovePeer(peer string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPeers := make([]string, 0)
	for _, p := range s.peers {
		if p != peer {
			newPeers = append(newPeers, p)
		}
	}
	s.peers = newPeers
}

func (s *RaftServer) AddPeer(peer string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.peers {
		if p == peer {
			return
		}
	}

	s.peers = append(s.peers, peer)
	s.nextIndex[peer] = int64(len(s.Log))
}
