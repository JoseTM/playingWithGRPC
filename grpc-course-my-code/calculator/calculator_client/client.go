package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close() // with defer this line executes at the end of program

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	doServerStream(c)

	doClientStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.CalculatorRequest{
		Operators: &calculatorpb.Operators{
			OperatorA: 4,
			OperatorB: 5,
		},
	}
	result, err := c.Add(context.Background(), &req)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	fmt.Printf("sum of %v + %v = %v \n", req.GetOperators().GetOperatorA(), req.GetOperators().GetOperatorB(),
		result.GetResult())
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.PrimeDecompositionRequest{
		Operator: 120,
	}
	result, err := c.PrimeDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	for {
		msg, err := result.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something was wrong with response")
		}
		fmt.Printf("factor: %s \n", msg.GetResult())
	}

}

func doClientStream(c calculatorpb.CalculatorServiceClient) {

	fmt.Printf("sending values to calculate average: \n")
	requestList := []*calculatorpb.AverageRequest{
		{
			Operator: 100,
		},
		{
			Operator: 50,
		},
		{
			Operator: 25,
		},
	}
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	for _, req := range requestList {
		err := stream.Send(req)
		fmt.Printf("value: %d \n", req.GetOperator())
		if err != nil {
			log.Fatalf("something was wrong with sending client stream")
		}
		time.Sleep(1000 * time.Millisecond)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("something was wrong with response")
	}
	fmt.Printf("average: %f \n", msg.GetResult())

}
