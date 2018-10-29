package main

import (
	"encoding/json"
	"log"
	"time"
)

// Node represents a node in Raft protocol
type Node struct {
	currentTerm   int64
	lastHeartbeat int64
	peers         []string
	role          string
	nodeID        string
	votedFor      string
	votesGranted  int64

	electionTicker *time.Ticker
}

func (n *Node) Init() {
	ps := []string{}
	for _, p := range n.peers {
		if p != n.nodeID {
			ps = append(ps, p)
		}
	}
	n.peers = ps
}

// MarshalJSON implements json.Marshaler.
func (n *Node) MarshalJSON() ([]byte, error) {
	var data = struct {
		CurrentTerm   int64    `json:"currentTerm"`
		LastHeartbeat int64    `json:"lastHeartbeat"`
		Peers         []string `json:"peers"`
		Role          string   `json:"role"`
		NodeID        string   `json:"nodeId"`
		VotedFor      string   `json:"votedFor"`
		VotesGranted  int64    `json:"votesGranted"`
	}{
		CurrentTerm:   n.currentTerm,
		LastHeartbeat: n.lastHeartbeat,
		Peers:         []string{},
		Role:          n.role,
		NodeID:        n.nodeID,
		VotedFor:      n.votedFor,
		VotesGranted:  n.votesGranted,
	}
	return json.Marshal(data)
}

func (n *Node) start() {
	log.Printf("Start [%v]", n.nodeID)
	n.electionTicker = time.NewTicker(10 * time.Millisecond)
	n.role = "follower"

	for {
		<-n.electionTicker.C

		switch n.role {
		case "follower":
			log.Printf("Follower [%v]", n.nodeID)
			n.startElection()

		case "leader":
			log.Printf("Leader [%v]", n.nodeID)
			n.updateFollowers()
		}
	}
}

func (n *Node) startElection() {
	n.role = "candidate"
	n.currentTerm++

	electionTerm := n.currentTerm
	votes := 1

	for _, p := range n.peers {
		p := p

		newTerm, isVoteGranted, err := n.requestVote(p, electionTerm, n.nodeID)
		if err != nil {
			log.Printf("RequestVote: %v", err)
		}

		log.Printf("RequestVote[%v]: %v %v", p, isVoteGranted, newTerm)

		if isVoteGranted {
			votes++
		}

		if newTerm > electionTerm {
			n.role = "follower"
			n.electionTicker.Stop()
			n.electionTicker = time.NewTicker(3500 * time.Millisecond)
			return
		}
	}

	if 2*votes > len(n.peers) {
		log.Printf("[%v] is a LEADER", n.nodeID)
		n.role = "leader"
		n.electionTicker.Stop()
		n.electionTicker = time.NewTicker(1500 * time.Millisecond)
		n.updateFollowers()
	} else {
		log.Printf("[%v] is a FOLLOWER", n.nodeID)
		n.role = "follower"
		n.electionTicker.Stop()
		n.electionTicker = time.NewTicker(5000 * time.Millisecond)
	}
}

func (n *Node) updateFollowers() {
	n.currentTerm++

	updateTerm := n.currentTerm

	for _, p := range n.peers {
		p := p

		newTerm, ok, err := n.appendEntries(p, updateTerm)
		if err != nil {
			log.Printf("AppendEntries: %v", err)
		}

		log.Printf("AppendEntries[%v]: %v %v", p, ok, newTerm)

		if !ok {
		}

		if newTerm > updateTerm {
			log.Printf("[%v] is a FOLLOWER", n.nodeID)
			n.role = "follower"
			n.electionTicker.Stop()
			n.electionTicker = time.NewTicker(5000 * time.Millisecond)
			return
		}
	}
}
