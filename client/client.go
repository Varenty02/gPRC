package main

import (
	"context"
	"gRPCDemo/proto/pb/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("err while dial %v", err)
	}
	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	//callSum(client)
	//callPND(client)
	//callAverage(client)
	//callMax(client)
	//callSquare(client, -64)
	callSumDeadline(client, 1*time.Second)
	callSumDeadline(client, 5*time.Second)
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
func callSumDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("calling sum with deadline exceed")
			}
		} else {
			log.Fatalf("call sum with deadline api err %v", err)
		}
		return
	}
	log.Println("sum api response %v", resp.GetResult())
}
func callPND(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(),
		&calculatorpb.PNDRequest{
			Number: 120,
		})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}
	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return
		}
		log.Println("prime number %v", resp.GetResult())
	}
}
func callAverage(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatal("call average err %v", err)
	}
	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{Num: 5},
		calculatorpb.AverageRequest{Num: 10.2},
		calculatorpb.AverageRequest{Num: 2.1},
		calculatorpb.AverageRequest{Num: 3.5},
		calculatorpb.AverageRequest{Num: 4},
		calculatorpb.AverageRequest{Num: 8.5},
		calculatorpb.AverageRequest{Num: 7},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("send average request err %v", err)
		}

	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receives average request err %v", err)
	}
	log.Printf("average response %+v", resp)
}
func callMax(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatal("call max err %v", err)
	}
	wait := make(chan struct{})
	go func() {
		listReq := []calculatorpb.MaxRequest{
			calculatorpb.MaxRequest{Num: 5},
			calculatorpb.MaxRequest{Num: 10.2},
			calculatorpb.MaxRequest{Num: 2.1},
			calculatorpb.MaxRequest{Num: 3.5},
			calculatorpb.MaxRequest{Num: 4},
			calculatorpb.MaxRequest{Num: 8.5},
			calculatorpb.MaxRequest{Num: 7},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send find max request err %v", err)
			}

		}
		stream.CloseSend()
	}()
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {

				break
			}
			if err != nil {
				log.Fatalf("recv find max err %v", err)
				break
			}
			log.Println("max %v", resp.GetResult())
		}
		close(wait)
	}()
	<-wait
}
func callSquare(c calculatorpb.CalculatorServiceClient, num float32) {
	resp, err := c.Square(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})
	if err != nil {
		log.Printf("call square api err %v\n", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err msg:%v", errStatus.Message())
			log.Printf("err status:%v", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Println("invalid argument num")
				return
			}
		}
	}
	log.Println("square api response %v", resp.GetResult())
}
