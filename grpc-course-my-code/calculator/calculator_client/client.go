package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
