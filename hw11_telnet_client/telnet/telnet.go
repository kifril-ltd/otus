package telnet

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrConnectionNotOpen = errors.New("connection is not open")

type Client struct {
	address              string
	timeout              time.Duration
	connectionRetryCount int
	sendRetryCount       int
	in                   io.Reader
	out                  io.Writer
	connection           net.Conn
}

type ClientOptionFunc func(*Client)

func WithConnectionRetry(count int) ClientOptionFunc {
	return func(client *Client) {
		if count == 0 {
			client.connectionRetryCount = 1
		} else {
			client.connectionRetryCount = count
		}
	}
}

func WithSendRetry(count int) ClientOptionFunc {
	return func(client *Client) {
		if count == 0 {
			client.sendRetryCount = 1
		} else {
			client.sendRetryCount = count
		}
	}
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
	options ...ClientOptionFunc,
) *Client {
	client := &Client{
		address:              address,
		timeout:              timeout,
		in:                   in,
		out:                  out,
		connectionRetryCount: 1,
		sendRetryCount:       1,
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

func (c *Client) Connect() error {
	var err error

	for i := 0; i < c.connectionRetryCount; i++ {
		conn, err := net.DialTimeout("tcp", c.address, c.timeout)
		if err == nil {
			c.connection = conn
			return nil
		}
	}

	return err
}

func (c *Client) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}

	return nil
}

func (c *Client) Send() error {
	if c.connection == nil {
		return ErrConnectionNotOpen
	}

	var err error
	for i := 0; i < c.sendRetryCount; i++ {
		_, err = io.Copy(c.connection, c.in)
		if err == nil {
			return nil
		}
	}

	return err
}

func (c *Client) Receive() error {
	if c.connection == nil {
		return ErrConnectionNotOpen
	}

	_, err := io.Copy(c.out, c.connection)

	return err
}
