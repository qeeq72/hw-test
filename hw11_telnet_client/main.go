package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration
	var host string
	var port string

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("\nUndefined arguments!\nPlease, use pattern:\n\tgo-telnet [--timeout=10] host port")
	}
	host = args[0]
	port = args[1]

	address := net.JoinHostPort(host, port)
	client, err := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect()
	if err != nil {
		log.Print("connection: ", err)
		return
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "....Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errCh := make(chan error, 1)
	defer func() {
		err := <-errCh
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer cancel()
		err := client.Send()
		if err == nil {
			errCh <- errors.New("...EOF")
		} else {
			errCh <- err
		}
	}()

	go func() {
		defer cancel()
		err := client.Receive()
		if err == nil {
			errCh <- errors.New("...Connection was closed by peer")
		} else {
			errCh <- err
		}
	}()

	<-ctx.Done()
}
