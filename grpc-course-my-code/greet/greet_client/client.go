package main

import (
	"context"
	"fmt"
	"log"

	"../greetpb"
	"google.golang.org/grpc"
)

type client struct{}

func main() {
	fmt.Println("Hello world")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close() // with defer this line executes at the end of program

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("create4d client: %f", c)

	myClient := client{}
	myClient.doUnary(c)

}

func (*client) doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jose",
			LastName:  "Trujillo",
		},
	}

	result, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not connect with error %v", err)
	}

	fmt.Printf("The result was %s \n", result.GetResult())
}
