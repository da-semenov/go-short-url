package app

import (
	"fmt"
	"net/http"
)

func RunApp() {
	h := ShortURLHandler()
	http.HandleFunc("/", h.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("service can't be started")
		fmt.Println(err)
	}
}
