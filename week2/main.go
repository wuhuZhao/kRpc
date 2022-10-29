package main

import (
	"fmt"
	"time"
	"week2/client"
	"week2/server"
)

func call() {
	time.Sleep(5 * time.Second)
	cli := client.NewKrpcClient()
	cli.Call(map[string]string{
		"context": "I am krpcClient",
		"date":    time.UnixDate,
	})
}

func main() {
	srv := server.NewKrpcServer(server.NewKrpcOptionWithDefaultMode("127.0.0.1", "10011", "tcp"))
	err := srv.Serve()
	go call()
	select {
	case er := <-err:
		{
			fmt.Printf("[krpcServer] connect error %+v", er)
		}
	}

}
