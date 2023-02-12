package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/thelazylemur/cacheengine/cache"
	"github.com/thelazylemur/cacheengine/client"
)

func main(){
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

	go func (){
		time.Sleep(time.Second * 2)
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 1; i++ {
			sendCommand(client)
			time.Sleep(time.Millisecond * 100)
		}

		client.Close()
		time.Sleep(time.Second * 1)
	}()

	server := NewServer(opts, cache.New())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func sendCommand(client *client.Client){
	_, err := client.Set(context.Background(), []byte("Dan"), []byte("Goat"), 0)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	client.Get(context.Background(), []byte("Dan"))
	if err != nil {
		log.Fatal(err)
	}
}
