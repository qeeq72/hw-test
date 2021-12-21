package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	DefaultHost    = "localhost"
	DefaultPort    = "8080"
	DefaultTimeout = 5 * time.Second
)

func main() {
	var timeout time.Duration
	var host string
	var port string

	flag.DurationVar(&timeout, "timeout", DefaultTimeout, "connection timeout")
	flag.Parse()

	args := flag.Args()
	switch len(args) {
	case 0:
		host = DefaultHost
		port = DefaultPort
	case 1:
		host = args[0]
		port = DefaultPort
	default:
		host = args[0]
		port = args[1]
	}

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

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Println("send: ", err)
		} else {
			fmt.Fprintln(os.Stderr, "...EOF")
		}
		cancel()
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Println("receive: ", err)
		} else {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		}
		cancel()
	}()

	<-ctx.Done()
}
