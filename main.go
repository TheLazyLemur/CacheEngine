package main

import (
	"flag"
	"log"
	"net"

	"github.com/thelazylemur/cacheengine/cache"
)

func main(){
	//testfollower()
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader node")
	)

	flag.Parse()

	opts := ServerOpts {
		ListenAddr: *listenAddr,
		IsLeader: len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	server := NewServer(opts, cache.New())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func testfollower(){
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

        _, err = conn.Write([]byte("SET Foo Bar 40000000"))
	if err != nil {
		log.Fatal(err)
	}

	select{}
}
