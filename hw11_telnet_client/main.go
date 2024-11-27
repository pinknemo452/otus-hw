package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
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

	buf := bytes.Buffer{}
	client := NewTelnetClient(addr, *timeout, io.NopCloser(&buf), os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buf.WriteString(scanner.Text() + "\n")
			err := client.Send()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to send: %v", err)
			}
		}
		if errors.Is(scanner.Err(), io.EOF) {
			println("...EOF")
			client.Close()
			os.Exit(0)
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
}
