package main

import (
	"fmt"
	"log"

	"../greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello world")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close() // with defer this line executes at the end of program

	c := greetpb.NewGreetServiceClient(cc)

	fmt.Printf("created client: %f", c)
}
