package main

import (
	"flag"
	"log"

	"context"
	"fmt"
	"os"

	"github.com/insighted4/insighted-go/examples/github/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	serverAddr := flag.String("server_addr", "127.0.0.1:8081", "The server address in the format of host:port")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"username\"\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("unable to close gRPC connection: %v", err)
		}
	}()

	client := api.NewGithubProxyClient(conn)

	req := &api.GetUserRequest{Name: flag.Arg(0)}
	res, err := client.GetUser(context.Background(), req)
	if err != nil {
		grpclog.Fatalf("could not get user: %v", err)
	}

	fmt.Println(res)
}
