package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	conn net.Conn
	err  error

	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	logger io.Writer
}

var ErrNoConnection = errors.New("no connection established")

func (c *Client) Connect() error {
	if c.conn != nil {
		return nil
	}

	c.conn, c.err = net.DialTimeout("tcp", c.address, c.timeout)

	return c.done("connect", "Connected to "+c.address)
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	c.err = c.conn.Close()
	c.conn = nil

	return c.done("close", "")
}

func (c *Client) Send() error {
	if c.conn == nil {
		return ErrNoConnection
	}

	_, c.err = io.Copy(c.conn, c.in)

	return c.done("send", "EOF")
}

func (c *Client) Receive() error {
	if c.conn == nil {
		return ErrNoConnection
	}

	_, c.err = io.Copy(c.out, c.conn)

	return c.done("receive", "Connection was closed by peer")
}

func (c *Client) done(msg string, log string) error {
	if c.err != nil {
		err := fmt.Errorf("%s: %w", msg, c.err)
		c.err = nil

		return err
	}

	if log != "" && c.logger != nil {
		fmt.Fprintf(c.logger, "...%s\n", log)
	}

	return nil
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
	logger io.Writer,
) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		logger:  logger,
	}
}
