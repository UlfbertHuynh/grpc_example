package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"math"
	pb "myGrpcExample/gen/mymsg"
	"net"
)

type server struct{
	pb.UnimplementedMenuOrderServer
}

var(
	promoMap = map[int64]float64{ 11028: 0.2, 45067: 0.1, 1212: 0.25}
	dishInfoMap = map[string]*DishInfo{
		"chicken_rice" : {Price: 3.0, PrepareTime: 2, Stock: 20},
		"beef_stew": {Price: 11.5, PrepareTime: 20, Stock: 30},
		"lamb_steak" : {Price: 9, PrepareTime: 20, Stock: 24},
		"salmon_salad" : {Price: 7.5, PrepareTime: 10, Stock: 15},
	}
)

type DishInfo struct{
	Price       float64
	PrepareTime int64 // in minutes
	Stock       int64
}

func (s *server) OrderChickenRice(ctx context.Context, request *pb.OrderRequest) (*pb.OrderResponse, error){
	log.Println("Received order for Chicken Rice")
	dishInfo, _ := dishInfoMap["chicken_rice"]

	if err := CheckStock(dishInfo, request.NumServing); err != nil{
		return nil, err
	}

	log.Printf("stock for remain: %d", dishInfoMap["chicken_rice"].Stock)
	return &pb.OrderResponse{
		Status: true,
		WaitingTime: dishInfo.PrepareTime * request.NumServing,
		TotalCost: CalculatePriceAfterPromo(dishInfo.Price, request.NumServing, request.PromoCode),
	}, nil
}

func (s * server) OrderBeefStew(ctx context.Context, request *pb.OrderRequest) (*pb.OrderResponse, error){
	log.Println("Received order for Beef Stew")
	dishInfo, _ := dishInfoMap["beef_stew"]

	if err := CheckStock(dishInfo, request.NumServing); err != nil{
		return nil, err
	}

	log.Printf("stock for remain: %d", dishInfoMap["beef_stew"].Stock)
	return &pb.OrderResponse{
		Status: true,
		WaitingTime: dishInfo.PrepareTime + int64(math.Ceil(float64(dishInfo.PrepareTime* (request.NumServing - 1)) * 0.1)),
		TotalCost: CalculatePriceAfterPromo(dishInfo.Price, request.NumServing, request.PromoCode),
	}, nil
}

func (s *server) OrderLambSteak(ctx context.Context, request *pb.OrderRequest) (*pb.OrderResponse, error){
	log.Println("Received order for Lamb Steak")
	dishInfo, _ := dishInfoMap["lamb_steak"]

	if err := CheckStock(dishInfo, request.NumServing); err != nil{
		return nil, err
	}

	log.Printf("stock for remain: %d", dishInfoMap["lamb_steak"].Stock)
	return &pb.OrderResponse{
		Status: true,
		WaitingTime: dishInfo.PrepareTime + int64(math.Ceil(float64(dishInfo.PrepareTime* (request.NumServing - 1)) * 0.15)),
		TotalCost: CalculatePriceAfterPromo(dishInfo.Price, request.NumServing, request.PromoCode),
	}, nil
}

func (s * server) OrderSalmonSalad(ctx context.Context, request *pb.OrderRequest) (*pb.OrderResponse, error){
	log.Println("Received order for Salmon Salad")
	dishInfo, _ := dishInfoMap["salmon_salad"]

	if err := CheckStock(dishInfo, request.NumServing); err != nil{
		return nil, err
	}

	log.Printf("stock for remain: %d", dishInfoMap["salmon_salad"].Stock)
	return &pb.OrderResponse{
		Status: true,
		WaitingTime: dishInfo.PrepareTime + int64(math.Ceil(float64(dishInfo.PrepareTime* (request.NumServing - 1)) * 0.05)),
		TotalCost: CalculatePriceAfterPromo(dishInfo.Price, request.NumServing, request.PromoCode),
	}, nil
}

func CalculatePriceAfterPromo(dishPrice float64, nServing int64, promoCode int64) float64{
	return dishPrice * float64(nServing) * (1 - promoMap[promoCode])
}

func CheckStock(dishInfo *DishInfo, nServing int64) error{
	if dishInfo.Stock < nServing{
		return errors.New("sorry stock run out")
	}
	dishInfo.Stock -= nServing
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMenuOrderServer(s, &server{})
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}