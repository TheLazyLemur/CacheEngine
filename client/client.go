package client

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/TheLazyLemur/cacheengine/protocol"
)

type Options struct {
	threadSafe bool
}

func NewOptions(safe bool) *Options {
	return &Options{
		threadSafe: safe,
	}
}

type Client struct {
	conn net.Conn
	lock sync.Mutex
	Options
}

func New(url string, opt Options) (*Client, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:    conn,
		Options: opt,
	}

	if opt.threadSafe {
		c.lock = sync.Mutex{}
	}

	return c, nil
}

func NewFromConn(conn net.Conn, opt Options) (*Client, error) {
	c := &Client{
		conn:    conn,
		Options: opt,
	}

	if opt.threadSafe {
		c.lock = sync.Mutex{}
	}

	return c, nil
}

func (c *Client) Set(_ context.Context, key, value []byte, ttl int) error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	cmd := &protocol.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := protocol.ParseSetReponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != protocol.StatusOK {
		return fmt.Errorf("server response with a non ok status: %s", resp.Status)
	}

	return nil
}

func (c *Client) Get(_ context.Context, key []byte) ([]byte, error) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	cmd := &protocol.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := protocol.ParseGetReponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status == protocol.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key %s", key)
	}

	if resp.Status != protocol.StatusOK {
		return nil, fmt.Errorf("server response with a non ok status: %s", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Delete(_ context.Context, key []byte) error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	cmd := &protocol.CommandDel{Key: key}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := protocol.ParseDelReponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != protocol.StatusOK {
		return fmt.Errorf("server response with a non ok status: %s", resp.Status)
	}

	return nil
}

func (c *Client) All(_ context.Context) ([][]byte, error) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	cmd := &protocol.CommandAll{}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := protocol.ParseAllResponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status != protocol.StatusOK {
		return nil, fmt.Errorf("server response with a non ok status: %s", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Close() error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	return c.conn.Close()
}
