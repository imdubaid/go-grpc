package main

import (
	"context"
	"log"
	"net"

	pb "github.com/imdubaid/grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(req *pb.MenuRequest, stream pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{
		{Id: "1", Name: "Espresso"},
		{Id: "2", Name: "Cappuccino"},
		{Id: "3", Name: "Latte"},
	}

	for i := range items {
		if err := stream.Send(&pb.Menu{Items: items[0 : i+1]}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{Id: "123"}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{OrderId: "123", Status: "In Progress"}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server is running on port 50051")

	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeShopServer(grpcServer, &server{})

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
