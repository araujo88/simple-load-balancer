package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
		fmt.Println("Hello, world!")
	})

	port := os.Args[1]

	fmt.Println("Server running on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
