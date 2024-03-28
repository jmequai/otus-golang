package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stimeout := flag.String("timeout", "10s", "enter connection timeout")

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatalln("enter only host and port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	timeout, errTimeout := time.ParseDuration(*stimeout)
	if errTimeout != nil {
		log.Fatalln("enter valid timeout")
	}

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout, os.Stderr)

	errConnect := client.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
	}

	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		defer cancel()
		if err := client.Send(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-ctx.Done()
}
