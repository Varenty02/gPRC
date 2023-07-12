package main

import (
	"context"
	"gRPCDemo/proto/pb/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called...")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}
func main() {
	log.Println()
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatal("err while create listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("err while serve %v", err)
	}
}
