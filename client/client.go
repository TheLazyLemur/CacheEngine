package client

import (
	"context"
	"net"

	"github.com/thelazylemur/cacheengine/protocol"
)

type Options struct {
}

type Client struct {
	conn net.Conn	
}

func New(url string, opt Options) (*Client, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn: conn,
	}

	return	c, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) (any, error) {
	cmd := &protocol.CommandSet{
		Key: key,
		Value: value,
		TTL: ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) Get(ctx context.Context, key []byte) (any, error) {
	cmd := &protocol.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
