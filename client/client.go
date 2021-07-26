package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	pb "myGrpcExample/gen/mymsg"
	"time"
)

const(
	serverAddress = ":8080"
)

type OrderFunc func(ctx context.Context, in *pb.OrderRequest, opts ...grpc.CallOption) (*pb.OrderResponse, error)

var (
	customerId = flag.Int("i", 0, "customer ID")
	promoCode = flag.Int("p", 0, "promotion code if any")
	dishName = flag.String("name", "", "name of the dish")
	nServing = flag.Int("n", 1, "number of serving for the dish")
	desc = flag.String("d", "", "extra description if any")
)



func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMenuOrderClient(conn)
	mapOrder := map[string]OrderFunc{
		"chicken_rice" : c.OrderChickenRice,
		"beef_stew"  : c.OrderBeefStew,
		"lamb_steak" : c.OrderLambSteak,
		"salmon_salad" : c.OrderSalmonSalad,
	}

	orderFunc, ok := mapOrder[*dishName]
	if *customerId == 0 || !ok{
		log.Fatal("Invalid customerId or dish name")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := orderFunc(ctx, &pb.OrderRequest{
		CustomerID:  int64(*customerId),
		PromoCode:   int64(*promoCode),
		Description: *desc,
		NumServing:  int64(*nServing),
	})
	if err != nil {
		log.Fatalf("could not accept order: %v", err)
	}
	log.Printf("order accepted, expected waiting time: %d, total cost: %f", r.WaitingTime, r.TotalCost)
}