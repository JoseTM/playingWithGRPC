package main

import (
	"context"
	"fmt"
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

	req := calculatorpb.CalculatorRequest{
		Operators: &calculatorpb.Operators{
			OperatorA: 4,
			OperatorB: 5,
		},
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)

	result, err := c.Add(context.Background(), &req)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	fmt.Printf("sum of %v + %v = %v \n", req.GetOperators().GetOperatorA(), req.GetOperators().GetOperatorB(),
		result.GetResult())
}
