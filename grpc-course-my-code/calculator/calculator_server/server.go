package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

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

func (*server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	primeNumber := req.GetOperator()
	fmt.Printf("prime decomposition in progress to %v \n", primeNumber)

	i := int64(2)
	for primeNumber > 1 {
		if primeNumber%i == 0 {
			result := calculatorpb.PrimeDecompositionResult{
				Result: strconv.FormatInt(i, 10),
			}
			stream.Send(&result)
			primeNumber = primeNumber / i
		} else {
			i++
		}

	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Printf("Greet function was invoked with a client streaming request to calculagte average")
	sum := int64(0)
	qty := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := float32(sum) / float32(qty)
			return stream.SendAndClose(&calculatorpb.AverageResult{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("somthing was wrong: %v", err)
		}
		sum += req.GetOperator()
		qty++
	}

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
