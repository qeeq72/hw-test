package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	errNilReader = errors.New("nil reader")
	errNilWriter = errors.New("nil writer")
)

type TelnetClientErrors []error

func (tce TelnetClientErrors) Error() string {
	n := len(tce)
	if n < 1 {
		return ""
	}
	s := "telnet client fail:\n"
	for i := 0; i < n; i++ {
		s += fmt.Sprintf("%s\n", tce[i])
	}
	return s
}

type TelnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	input   io.ReadCloser
	output  io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) (*TelnetClient, error) {
	var errs TelnetClientErrors
	if in == nil {
		errs = append(errs, errNilReader)
	}
	if out == nil {
		errs = append(errs, errNilWriter)
	}
	if len(errs) > 0 {
		return nil, errs
	}

	return &TelnetClient{
		address: address,
		timeout: timeout,
		input:   in,
		output:  out,
	}, nil
}

func (c *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TelnetClient) Close() error {
	return c.conn.Close()
}

func (c *TelnetClient) Send() error {
	_, err := io.Copy(c.conn, c.input)
	return err
}

func (c *TelnetClient) Receive() error {
	_, err := io.Copy(c.output, c.conn)
	return err
}
