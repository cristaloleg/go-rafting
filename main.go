package main

import (
	"flag"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

func main() {
	port := flag.String("port", "3133", "")
	flag.Parse()

	node := Node{
		role:   "über-leader",
		nodeID: "127.0.0.1:" + *port,
		// nodeID: "172.18.200.93:" + *port,
		peers: []string{
			// "172.18.203.132:5000",
			// "172.18.203.61:6666",
			// "172.18.201.15:8080",
			// "172.18.203.23:8080",
			// "172.18.200.141:8001",
			// "172.18.202.201:9999",
			// "172.18.201.70:9999",

			"127.0.0.1:3101",
			"127.0.0.1:3102",
			"127.0.0.1:3103",
		},
	}

	s := Peer{
		node: node,
	}
	s.Start(*port)
}
