package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Println("main : started")
	defer log.Println("main : completed")

	//start api
	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	//if err := http.ListenAndServe("localhost:8000", h); err != nil {
	//	log.Fatal(err)
	//}
	go func() {
		log.Println("API listening on address:", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	//Make a channel to listen for interrupt or terminate signal from OS.
	//Use the buffered channel because the signal package requeres it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	//===========================================
	select {
	case err := <-serverErrors:
		log.Fatalf("Error : listening and serving : %s", err)
	case <-shutdown:
		log.Println("main : start shutdown")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Println("main : gracefull shutdown")
			err = api.Close()
		}
		if err != nil {
			log.Println("main : could not stop server grasefully")
		}
	}
}

func Echo(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(1000)

	fmt.Println("starting ", id)

	time.Sleep(3 * time.Second)

	fmt.Fprintln(w, "You asked for:", r.Method, r.URL.Path)

	fmt.Println("ending ", id)
}
