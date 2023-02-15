package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/TheLazyLemur/cacheengine/cache"
	"github.com/TheLazyLemur/cacheengine/client"
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
		testClient()

	}()

	server := NewServer(opts, cache.New())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func testClient(){
	for i := 0; i <= 10; i++ {
		go func(i int){
			var (
				key = []byte(fmt.Sprintf("key_%d", i))
				val = []byte(fmt.Sprintf("val_%d", i))
			)
			client, err := client.New(":3000", client.Options{})

			if err != nil {
				log.Fatal(err)
			}

			err = client.Set(context.Background(), key, val, 0)
			if err != nil {
				log.Fatal(err)
			}

			err = client.Delete(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			
			resp, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(resp))

			client.Close()
		}(i)
	}
}
