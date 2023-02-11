package main
//
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net"
//
// 	"github.com/thelazylemur/cacheengine/cache"
// )
//
// type ServerOpts struct {
// 	ListenAddr string
// 	IsLeader bool
// 	LeaderAddr string
// }
//
// type Server struct {
// 	ServerOpts
// 	followers map[net.Conn]struct{}
// 	cacher cache.Cacher
// }
//
// func NewServer(opts ServerOpts, c cache.Cacher) *Server {
// 	return &Server{
// 		ServerOpts: opts,
// 		cacher: c,
// 		//TODO: only allocate this as the leader
// 		followers: make(map[net.Conn]struct{}),
// 	}
// }
//
// func (s *Server) Start() error {
// 	fmt.Println(s.ListenAddr)
// 	ln, err := net.Listen("tcp", s.ListenAddr)
// 	if err != nil {
// 		return fmt.Errorf("listen error: %s", err)
// 	}
//
// 	log.Printf("server starting on port [%s]\n", s.ListenAddr)
//
// 	if !s.IsLeader {
// 		go func() {
// 			conn, err := net.Dial("tcp", s.LeaderAddr)
// 			fmt.Println("Connected with leader:", s.LeaderAddr)
// 			if err != nil {
// 				log.Fatal(err.Error())
// 			}
//
// 			s.handleConn(conn)
// 		}()
// 	}
//
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			log.Printf("accept error: %s\n", err)
// 			continue
// 		}
//
// 		go s.handleConn(conn)
// 	}
// }
//
// func (s *Server) handleConn(conn net.Conn) {
// 	defer func() {
// 		if err := conn.Close(); err != nil {
// 			log.Printf("close error: %s\n", err)
// 		}	
// 	}()
//
// 	buf := make([]byte, 2048)
//
// 	if s.IsLeader {
// 		s.followers[conn] = struct{}{}
// 	}
//
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			log.Printf("read error: %s\n", err)
// 			break
// 		}
//
// 		go s.handleCommand(conn, buf[:n])
// 	}
// }
//
// func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
// 	msg, err := parseMessage(rawCmd)
// 	if err != nil {
// 		log.Println("error parsing command")
// 		return
// 	}
//
// 	switch msg.Cmd {
// 	case CMDSet:
// 		err = s.handleSetCommand(conn, msg)
// 	case CMDGet: 
// 		err = s.handleGetCommand(conn, msg)
// 	case CMDHas:
// 		err = s.handleHasCommand(conn, msg)
// 	case CMDDel:
// 		err = s.handleDelCommand(conn, msg)
// 	default:
// 		err = errors.New("unknown command")
// 	}
//
// 	if err != nil {
// 		log.Println(err.Error())
// 		s.writeError(conn, err)
// 	}else if msg.Cmd == CMDSet || msg.Cmd == CMDDel {
// 		go func() {
// 			if err := s.sendToFollowers(context.TODO(), msg); err != nil {
// 				log.Println("error sending to followers")
// 			}
// 		}()
// 	}
// }
//
// func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
// 	return s.cacher.Set(msg.Key, msg.Value, msg.Ttl);
// }
//
// func (s *Server) handleGetCommand(conn net.Conn, msg *Message) error {
// 	val, err := s.cacher.Get(msg.Key);
// 	if err != nil {
// 		return err
// 	}
//
// 	_, _ = conn.Write(val)
//
// 	return nil
// }
//
// func (s *Server) handleHasCommand(conn net.Conn, msg *Message) error {
// 	ok := s.cacher.Has(msg.Key);
// 	if ok {
// 		_, _ = conn.Write([]byte("ok"))
// 	}
//
// 	return nil
// }
//
// func (s *Server) handleDelCommand(conn net.Conn, msg *Message) error {
// 	s.cacher.Delete(msg.Key);
//
// 	return nil
// }
//
// func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
// 	for conn := range s.followers {
// 		_, err := conn.Write(msg.ToBytes())
// 		if err != nil {
// 			log.Println("error writing to follower", err.Error())
// 			continue
// 		}
// 	}
// 	return nil
// }
//
// func (s *Server) writeError(conn net.Conn, err error) {
// 	log.Println(err.Error())
// 	_, _ = conn.Write([]byte(err.Error()))
// }
