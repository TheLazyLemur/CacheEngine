package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TheLazyLemur/cacheengine/api"
	"github.com/TheLazyLemur/cacheengine/cache"
	"github.com/TheLazyLemur/cacheengine/client"
	"github.com/TheLazyLemur/cacheengine/server"
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

	opts := server.ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	apiOpts := api.ApiServerOpts{
		ListenAddr: *apiAddr,
	}

	c := cache.New()

	api := api.NewApiServer(apiOpts, c)
	go api.Run()

	server := server.NewServer(opts, c)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func testClient() {
	wg := &sync.WaitGroup{}

	client, err := client.New(":3000", *client.NewOptions(true))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := []byte("Foo" + fmt.Sprint(i))
			val := []byte("Bar" + fmt.Sprint(i))

			err := client.Set(context.Background(), key, val, 0)
			if err != nil {
				log.Fatal(err)
			}

			resp, err := client.Get(context.Background(), key)
			if err != nil {
				log.Println("key not found")
			} else {
				fmt.Println(string(resp))
			}

			keys, err := client.All(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			for _, k := range keys {
				fmt.Println("\t", string(k))
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("done")
}
