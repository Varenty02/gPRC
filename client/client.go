package main

import (
	"context"
	"gRPCDemo/proto/pb/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("err while dial %v", err)
	}
	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	callSum(client)
}
func callSum(c calculatorpb.CalculatorServiceClient) {
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		log.Fatal("call sum api error %v", err)
	}
	log.Println("sum api response %v", resp.GetResult())
}
