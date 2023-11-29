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

	"github.com/kifril-ltd/otus-hw/hw11_telnet_client/telnet"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func main() {
	var (
		timeout         time.Duration
		connectionRetry int
		sendRetry       int
	)
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connect")
	flag.IntVar(&connectionRetry, "connection-retry", 1, "number of connection retry")
	flag.IntVar(&sendRetry, "send-retry", 1, "number of send retry")

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatalln("Wrang args count ", args)
	}

	address := net.JoinHostPort(args[0], args[1])
	client := telnet.NewTelnetClient(
		address, timeout, os.Stdin, os.Stdout,
		telnet.WithConnectionRetry(connectionRetry),
		telnet.WithSendRetry(sendRetry),
	)

	if err := run(client); err != nil {
		log.Fatalln(err)
	}
}

func run(client TelnetClient) (err error) {
	if err = client.Connect(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	defer func() {
		if err := client.Close(); err != nil {
			cancel()
		}
	}()

	go func() {
		defer cancel()
		err = client.Send()
	}()

	go func() {
		defer cancel()
		err = client.Receive()
	}()

	<-ctx.Done()

	return nil
}
