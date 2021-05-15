package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../greetpb"

	"google.golang.org/grpc"
)

type client struct{}

func main() {
	// fmt.Println("Hello world")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close() // with defer this line executes at the end of program

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("create4d client: %f", c)

	myClient := client{}
	/* 	myClient.doUnary(c)
	   	myClient.doServerStream(c)
	   	myClient.doClientStream(c) */
	myClient.doBiStream(c)

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

func (*client) doServerStream(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetingManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jose",
			LastName:  "Trujillo",
		},
	}

	result, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not connect with error %v", err)
	}

	for {
		msg, err := result.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("somthing was wrong with response")
		}
		fmt.Printf("menssage stream: %s \n", msg.GetResult())
	}

}

func (*client) doClientStream(c greetpb.GreetServiceClient) {

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Could not connect with error %v", err)
	}

	reqlist := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jose",
				LastName:  "Trujillo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "María",
				LastName:  "Trujillo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lorena",
				LastName:  "Trujillo",
			},
		},
	}

	for _, req := range reqlist {
		err := stream.Send(req)
		fmt.Printf("sending menssage client stream: %s\n", req.GetGreeting().FirstName)
		if err != nil {
			log.Fatalf("somthing was wrong sending")
		}
		time.Sleep(1000 * time.Millisecond)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("somthing was wrong with response")
	}
	fmt.Printf("menssage client stream:\n%s", msg.GetResult())

}

func (*client) doBiStream(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Could not connect with error %v", err)
	}

	reqlist := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jose",
				LastName:  "Trujillo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "María",
				LastName:  "Trujillo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lorena",
				LastName:  "Trujillo",
			},
		},
	}

	waitChannel := make(chan struct{})

	go func() {
		for _, req := range reqlist {
			fmt.Printf("sending menssage client stream: %s\n", req.GetGreeting().FirstName)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("boom")
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("somthing was wrong with response")
				break
			}
			fmt.Printf("menssage stream: %s \n", msg.GetResult())
		}
		close(waitChannel)
	}()

	<-waitChannel

}
