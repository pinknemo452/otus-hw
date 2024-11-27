package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}
type Client struct {
	sync.Mutex
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {

	return &Client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *Client) Connect() error {
	tc.Lock()
	defer tc.Unlock()
	d := net.Dialer{Timeout: tc.timeout}
	conn, err := d.Dial("tcp", tc.addr)
	if err != nil {
		return err
	}
	tc.conn = conn

	println("...Connected to", tc.addr)
	return nil
}

func (tc *Client) Send() error {
	readed, err := io.ReadAll(tc.in)
	if err != nil {
		return err
	}
	log.Printf("send: readed %d bytes", len(readed))
	writen, err := tc.conn.Write([]byte(readed))
	if err != nil {
		return err
	}
	log.Printf("send: written %d bytes", writen)
	if len(readed) != writen {
		return errors.New("failed to send whole message")
	}
	return nil
}

func (tc *Client) Receive() error {
	buf := make([]byte, 32768)
	readed, err := tc.conn.Read(buf)
	if err != nil {
		return err
	}
	log.Printf("receive: read %d bytes", readed)
	buf = buf[:readed]
	written, err := tc.out.Write(buf)
	if err != nil {
		return err
	}
	log.Printf("receive: written %d bytes", written)
	if readed != written {
		return errors.New("failed to read whole message")
	}
	return nil
}

func (tc *Client) Close() error {
	tc.Lock()
	defer tc.Unlock()

	return tc.conn.Close()

}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
