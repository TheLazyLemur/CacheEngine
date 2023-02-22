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

func main() {
	var (
		apiAddr    = flag.String("apiaddr", ":8080", "listen address of the api")
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader node")
	)
	flag.Parse()

	go func() {
		time.Sleep(time.Second * 2)
		testClient()

	}()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	apiOpts := ApiServerOpts{
		ListenAddr: *apiAddr,
	}

	c := cache.New()

	api := NewApiServer(apiOpts, c)
	go api.Run()

	server := NewServer(opts, c)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func testClient() {
	client, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 1; i++ {
		go func(i int) {
			key := []byte("Foo")
			key2 := []byte("Foo2")
			val := []byte("Bar")

			err = client.Set(context.Background(), key, val, 0)
			if err != nil {
				log.Fatal(err)
			}

			resp, err := client.Get(context.Background(), key)
			if err != nil {
				log.Println("key not found")
			} else {
				fmt.Println(string(resp))
			}

			err = client.Set(context.Background(), key2, val, 0)
			if err != nil {
				log.Fatal(err)
			}

			resp, err = client.Get(context.Background(), key2)
			if err != nil {
				log.Println("key not found")
			} else {
				fmt.Println(string(resp))
			}

			err = client.All(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			client.Close()
		}(i)
	}
}
