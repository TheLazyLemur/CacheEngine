package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/thelazylemur/cacheengine/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader bool
}

type Server struct {
	ServerOpts ServerOpts
	cacher cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cacher: c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ServerOpts.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ServerOpts.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close error: %s\n", err)
		}	
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error: %s\n", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		log.Println("error parsing command")
		return
	}

	switch msg.Cmd {
	case CMDSet:
		if err := s.handleSetCommand(conn, msg); err != nil {
			log.Println("something went wrong while handling the SET command: ", msg)
			return
		}
	case CMDGet: 
		value, err := s.handleGetCommand(conn, msg)
		if err != nil {
			log.Println("something went wrong while handling the GET command: ", msg)
		}
		fmt.Println(string(value))
	}

	go func() {
		if err := s.sendToFollowers(context.TODO(), msg); err != nil {
			log.Println("error sending to followers")
		}
	}()
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	return s.cacher.Set(msg.Key, msg.Value, msg.Ttl);
}

func (s *Server) handleGetCommand(conn net.Conn, msg *Message) ([]byte, error) {
	return s.cacher.Get(msg.Key);
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
