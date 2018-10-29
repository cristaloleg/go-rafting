package main

import (
	"log"
	"net/http"
)

// Peer represents a peer in Raft consensus.
type Peer struct {
	node Node
}

// Start activates a peer.
func (s *Peer) Start(port string) {
	s.node.Init()
	go s.node.start()

	http.HandleFunc("/raft/request-vote", s.requestVoteHandler)
	http.HandleFunc("/raft/append-entries", s.appendEntriesHandler)
	http.HandleFunc("/raft/state", s.stateHandler)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panic(err)
	}
}

func (s *Peer) requestVoteHandler(w http.ResponseWriter, r *http.Request) {
	dumpRequest(r)

	var req = struct {
		Term        int64  `json:"term"`
		CandidateID string `json:"candidateId"`
	}{}
	if err := decode(r.Body, &req); err != nil {
		http.Error(w, "meh, bad json", http.StatusBadRequest)
		return
	}
	var resp = struct {
		Term      int64 `json:"term"`
		IsGranted bool  `json:"voteGranted"`
	}{
		Term:      s.node.currentTerm,
		IsGranted: s.node.currentTerm < req.Term,
	}
	_ = encode(w, resp)
}

func (s *Peer) appendEntriesHandler(w http.ResponseWriter, r *http.Request) {
	dumpRequest(r)

	var req = struct {
		Term int64 `json:"term"`
	}{}

	if err := decode(r.Body, &req); err != nil {
		http.Error(w, "meh, bad json", http.StatusBadRequest)
		return
	}

	var resp = struct {
		Term      int64 `json:"term"`
		IsSuccess bool  `json:"success"`
	}{
		Term:      s.node.currentTerm,
		IsSuccess: s.node.currentTerm < req.Term,
	}
	_ = encode(w, resp)
}

func (s *Peer) stateHandler(w http.ResponseWriter, r *http.Request) {
	dumpRequest(r)

	var response = struct {
		CurrentTerm   int64    `json:"currentTerm"`
		LastHeartbeat int64    `json:"lastHeartbeat"`
		Peers         []string `json:"peers"`
		Role          string   `json:"role"`
		NodeID        string   `json:"nodeId"`
		VotedFor      string   `json:"votedFor"`
		VotesGranted  int64    `json:"votesGranted"`
	}{
		CurrentTerm:   s.node.currentTerm,
		LastHeartbeat: s.node.currentTerm,
		Peers:         s.node.peers,
		Role:          s.node.role,
		NodeID:        s.node.nodeID,
		VotedFor:      s.node.votedFor,
		VotesGranted:  s.node.votesGranted,
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	_ = encode(w, response)
}
