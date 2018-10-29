package main

import "bytes"

func (n *Node) requestVote(address string, term int64, candidateID string) (newTerm int64, voteGranted bool, err error) {
	buf := bytes.NewBufferString(`{
		"term": ` + itoa(term) + `, 
		"candidateId": "` + candidateID + `"
		}`)

	resp, err := httpClient.Post("http://"+address+"/raft/request-vote", "application/json", buf)
	if err != nil {
		return 0, false, err
	}

	var response = struct {
		Term      int64 `json:"term"`
		IsGranted bool  `json:"voteGranted"`
	}{}

	err = decode(resp.Body, &response)
	if err != nil {
		return 0, false, err
	}
	return response.Term, response.IsGranted, nil
}

func (n *Node) appendEntries(address string, term int64) (newTerm int64, ok bool, err error) {
	buf := bytes.NewBufferString(`{"term": ` + itoa(term) + `}`)

	resp, err := httpClient.Post("http://"+address+"/raft/append-entries", "application/json", buf)
	if err != nil {
		return 0, false, err
	}

	var response = struct {
		Term      int64 `json:"term"`
		IsSuccess bool  `json:"success"`
	}{}

	err = decode(resp.Body, &response)
	if err != nil {
		return 0, false, err
	}
	return response.Term, response.IsSuccess, nil
}

func (n *Node) getState(address string) error {
	resp, err := httpClient.Get(address + "/raft/state")
	if err != nil {
		return err
	}

	var response = struct {
		CurrentTerm   int64    `json:"currentTerm"`
		LastHeartbeat int64    `json:"lastHeartbeat"`
		Peers         []string `json:"peers"`
		Role          string   `json:"role"`
		NodeID        string   `json:"nodeId"`
		VotedFor      string   `json:"votedFor"`
		VotesGranted  int64    `json:"votesGranted"`
	}{}

	err = decode(resp.Body, &response)
	if err != nil {
		return err
	}

	return nil
}
