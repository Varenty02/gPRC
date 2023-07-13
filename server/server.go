package main

import (
	"context"
	"gRPCDemo/proto/pb/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"time"
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
func (*server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called...")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("context cancel")
			return nil, status.Errorf(codes.Canceled, "client cancel req")
		}
		time.Sleep(1 * time.Second)
	}
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}
func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	n := req.GetNumber()
	k := int32(2)
	for n > 1 {
		if n%k == 0 {
			n = n / k
			stream.Send(&calculatorpb.PNDResponse{Result: k})
		} else {
			k++
		}
	}
	return nil
}
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatal("err while Recv average %v", err)
		}
		total += req.GetNum()
		count++
	}
}
func (*server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	max := float32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal("err while Recv max %v", err)
			return err
		}
		num := req.GetNum()
		if num > max {
			max = num
		}
		stream.Send(&calculatorpb.MaxResponse{Result: max})
		if err != nil {
			return err
		}
	}
}
func (*server) Square(ctx context.Context, req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	num := req.GetNum()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Expect num>0,req num was %v", num)
	}
	return &calculatorpb.SquareResponse{Result: float32(math.Sqrt(float64(num)))}, nil
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
