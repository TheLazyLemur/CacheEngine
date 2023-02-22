package client

import (
	"context"
	"fmt"
	"net"

	"github.com/TheLazyLemur/cacheengine/protocol"
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

	return c, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) error {
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

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
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

func (c *Client) Delete(ctx context.Context, key []byte) error {
	cmd := &protocol.CommandDel{
		Key: key,
	}

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

func (c *Client) Join(ctx context.Context) error {
	cmd := &protocol.CommandJoin{}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) All(ctx context.Context) error {
	cmd := &protocol.CommandAll{}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := protocol.ParseAllReponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != protocol.StatusOK {
		return fmt.Errorf("server response with a non ok status: %s", resp.Status)
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
