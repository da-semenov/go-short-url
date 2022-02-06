package main

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {
	http.HandleFunc("/", HelloWorld)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("server can't be started")
		fmt.Println(err)
	}
}
