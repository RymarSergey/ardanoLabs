package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	h := http.HandlerFunc(Echo)

	log.Println("we liastening on port:8000")

	if err := http.ListenAndServe("localhost:8000", h); err != nil {
		log.Fatal(err)
	}
}

func Echo(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(1000)

	fmt.Println("starting ", id)

	time.Sleep(3 * time.Second)

	fmt.Fprintln(w, "You asked for:", r.Method, r.URL.Path)

	fmt.Println("ending ", id)
}
