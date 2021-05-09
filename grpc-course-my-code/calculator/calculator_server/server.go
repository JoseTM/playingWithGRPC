package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Add(ctx context.Context, in *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResult, error) {
	fmt.Printf("adding operators %v \n", in.GetOperators())
	sum := in.GetOperators().GetOperatorA() + in.GetOperators().GetOperatorB()

	result := calculatorpb.CalculatorResult{
		Result: int64(sum),
	}

	return &result, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
