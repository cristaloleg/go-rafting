poof:
	go run ./...

reqv:
	http --json POST localhost:31337/raft/request-vote @data-request-vote.json

appe:
	http --json POST localhost:31337/raft/append-entries @data-append-entries.json 

stat:
	http GET localhost:31337/raft/state