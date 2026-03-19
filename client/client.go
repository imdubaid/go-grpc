package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/imdubaid/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a new grpc client
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	// dont forget to close it
	defer conn.Close()

	// create a new coffee shop client from our generated code and pass in the connection created above
	c := pb.NewCoffeeShopClient(conn)

	// give us a context that we can cancel, but also a timeout just to illustrate a point
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// stream the menu
	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}

	// this next bit is some channel manipulation to let a go routine run with recieving messages from the stream.
	// there are other ways to handle this, but this is how I choose to handle it.
	done := make(chan bool)

	// We'll store the items here so that we can refer to them after streaming
	var items []*pb.Item

	// start a go routine that runs until we get an EOF from the stream.
	// We use this because our server sends us a partial menu as it builds up a menu in memory.
	// When we get an EOF, the stream is finished and we have the most up to date menu.
	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			// Store the last message's items for use later
			items = resp.Items
			log.Printf("Resp received: %v", resp.Items)
		}
	}()

	// Wait until that done channel has a message.
	<-done

	// Make a simple call to order all the items on the menu
	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	log.Printf("%v", receipt)

	// Make a simple call to get the order status.
	status, err := c.GetOrderStatus(ctx, receipt)
	log.Printf("%v", status)
}
