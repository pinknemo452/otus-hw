package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	port := os.Args[len(os.Args)-1]
	host := os.Args[len(os.Args)-2]

	addr := net.JoinHostPort(host, port)

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for {
			err := client.Send()
			if errors.Is(err, io.EOF) {
				println("...EOF")
				client.Close()
				os.Exit(0)
			}
		}
	}()
	go func() {
		for {
			err := client.Receive()
			if errors.Is(err, io.EOF) {
				println("...Connection was closed by peer")
				client.Close()
				os.Exit(0)
			}
		}
	}()

	<-ctx.Done()
	client.Close()
	stop()
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
