package main

import (
	"io"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8021", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)
	}
}
