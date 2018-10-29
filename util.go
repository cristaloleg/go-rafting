package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

var dump = false

func dumpRequest(r *http.Request) {
	if !dump {
		return
	}
	defer func() {
		log.Printf("\n-----------------------------\n\n")
	}()

	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("lol %v %v", err, string(b))
	}
	log.Printf("%v", string(b))
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}

func decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
