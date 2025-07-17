package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func main() {

	h := http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		header := rsp.Header()
		header.Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(rsp, "Hello world")		
	})
		
	err := cgi.Serve(h)

	if err != nil {
		fmt.Println(err)
	}
}
