package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/TheLazyLemur/cacheengine/cache"
	"github.com/TheLazyLemur/cacheengine/client"
	"github.com/TheLazyLemur/cacheengine/protocol"
)

type Opts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	Opts
	members map[*client.Client]struct{}
	cache   cache.Cacher
}

func NewServer(opts Opts, c cache.Cacher) *Server {
	return &Server{
		Opts:    opts,
		cache:   c,
		members: make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	fmt.Println(s.ListenAddr)
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			err := s.dialLeader()
			if err != nil {
				log.Println("failed to dial leader:", err)
			}
		}()
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("error dialing leader: [%s]", err)
	}

	log.Println("connected to leader:", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, protocol.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close error: %s\n", err)
		}
	}()

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
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		_ = s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		_ = s.handleGetCommand(conn, v)
	case *protocol.CommandDel:
		_ = s.handleDelCommand(conn, v)
	case *protocol.CommandJoin:
		_ = s.handleJoinCommand(conn, v)
	case *protocol.CommandAll:
		_ = s.handleAllCommand(conn)
	default:
		fmt.Println("default")
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) error {
	op := client.NewOptions(false)
	c, err := client.NewFromConn(conn, *op)
	if err != nil {
		return err
	}

	s.members[c] = struct{}{}
	fmt.Println("member just joined the cluster", conn.RemoteAddr())
	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	log.Printf("SET %s to %s with ttl of %d\n", cmd.Key, cmd.Value, cmd.TTL)

	if s.IsLeader {
		go func() {
			for member := range s.members {
				err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
				if err != nil {
					log.Panicln("failed to forward to memeber:", err)
				}
			}
		}()
	}

	resp := protocol.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, int64(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		_, err = conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	_, _ = conn.Write(resp.Bytes())
	fmt.Println("Finished responding")

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) error {
	// log.Printf("GET %s\n", cmd.Key)
	resp := protocol.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
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

func (s *Server) handleDelCommand(conn net.Conn, cmd *protocol.CommandDel) error {
	// log.Printf("DEL %s\n", cmd.Key)
	resp := protocol.ResponseDelete{}
	err := s.cache.Delete(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusError
		_, _ = conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK

	_, err = conn.Write(resp.Bytes())

	fmt.Println("DELETING")

	return err
}

func (s *Server) handleAllCommand(conn net.Conn) error {
	resp := protocol.ResponseAll{}
	x, _ := s.cache.All()
	ks := make([][]byte, 0)
	ks = append(ks, x...)
	resp.Status = protocol.StatusOK
	resp.Value = ks
	resp.AmountKeys = int32(len(x))
	_, err := conn.Write(resp.Bytes())
	return err
}
