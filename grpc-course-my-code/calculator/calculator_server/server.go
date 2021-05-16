package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"../calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
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
	fmt.Printf("Gunction was invoked with a client streaming request to calculagte average")
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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Function was invoked with a client streaming request to calculate Maximum\n")
	max := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("somthing was wrong: %v", err)
		}
		if max < req.GetCandidate() {
			max = req.GetCandidate()
		}
		errSend := stream.Send(&calculatorpb.FindMaximumResult{
			Result: max,
		})
		if errSend != nil {
			log.Fatalf("somthing was wrong: %v\n", errSend)
			return errSend
		}
	}
}

func (*server) SquareRoot(ctx context.Context, in *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("square operator %v \n", in.GetNumber())
	number := in.GetNumber()

	for i := 0; i < 4; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("cliente canceled request")
			return nil, status.Error(codes.Canceled, "cancelado por el cliente, tiempo excedido")
		}
		time.Sleep(1000 * time.Millisecond)
	}

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "The argument must be > 0")
	}

	result := calculatorpb.SquareRootResponse{
		SquareRoot: math.Sqrt(float64(number)),
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
