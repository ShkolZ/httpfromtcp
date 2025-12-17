package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShkolZ/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	mux, err := server.NewMux()
	if err != nil {
		log.Fatalln(err)
	}

	mux.Register("/", HelloWorldHandler)

	server, err := server.Serve(port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
