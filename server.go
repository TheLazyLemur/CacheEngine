package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/thelazylemur/cacheengine/cache"
	"github.com/thelazylemur/cacheengine/protocol"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	followers map[net.Conn]struct{}
	cacher cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cacher: c,
		//TODO: only allocate this as the leader
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	fmt.Println(s.ListenAddr)
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Println("Connected with leader:", s.LeaderAddr)
			if err != nil {
				log.Fatal(err.Error())
			}

			s.handleConn(conn)
		}()
	}

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

	// fmt.Printf("new connection from [%s]\n", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			log.Printf("parse error: %s, dropping conection\n", err)
			break
		}

		go s.handleCommand(conn, cmd)
	}

	// fmt.Println("connection closed")
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		_ = s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		_ = s.handleGetCommand(conn, v)
	case *protocol.CommandDel:
		_ = s.handleDelCommand(conn, v)
	}
}

func (s *Server) handleSetCommand (conn net.Conn, cmd *protocol.CommandSet) error {
	// log.Printf("SET %s to %s\n", cmd.Key, cmd.Value)

	resp := protocol.ResponseSet{}
	if err := s.cacher.Set(cmd.Key, cmd.Value, int64(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		_, _ = conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	_, _ = conn.Write(resp.Bytes())
	return nil
}

func (s *Server) handleGetCommand (conn net.Conn, cmd *protocol.CommandGet) error {
	// log.Printf("GET %s\n", cmd.Key)

	resp := protocol.ResponseGet{}
	value, err := s.cacher.Get(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusKeyNotFound
		_, _ = conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	resp.Value = value

	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleDelCommand (conn net.Conn, cmd *protocol.CommandDel) error {
	// log.Printf("DEL %s\n", cmd.Key)
	resp := protocol.ResponseDelete{}
	err := s.cacher.Delete(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusError
		_, _ = conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK

	_, err = conn.Write(resp.Bytes())

	return err
}
