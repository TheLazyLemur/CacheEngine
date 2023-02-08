package main

import (
	"log"
	"net"
	"time"

	"github.com/thelazylemur/cacheengine/cache"
)

func main(){
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}

	go test()

	server := NewServer(opts, cache.New())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func test(){
	time.Sleep(time.Second * 2)

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte("SET Foo Bar 25000000000000000"))

	time.Sleep(time.Second * 2)

	conn.Write([]byte("HAS Foo"))

	time.Sleep(time.Second * 2)

	conn.Write([]byte("GET Foo"))
}
