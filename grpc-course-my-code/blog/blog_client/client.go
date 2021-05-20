package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../blogpb"

	"google.golang.org/grpc"
)

type client struct{}

func (*client) createBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "jose",
			Title:    "titulo " + time.Now().String(),
			Content:  "content",
		},
	}

	result, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not create blog with error %v", err)
	}

	fmt.Printf("The result was %s \n", result.GetBlog())
}

func (*client) ReadBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.ReadBlogRequest{
		BlogId: "60a44e5360744a50e140efc9",
	}

	result, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not get blog with error %v", err)
	}

	fmt.Printf("The result was %s \n", result.GetBlog())
}

func (*client) UpdateBlog(c blogpb.BlogServiceClient) {
	blog := &blogpb.Blog{
		Id:       "60a44e5360744a50e140efc9",
		AuthorId: "Gilberto " + time.Now().String(),
		Title:    "otro título",
		Content:  time.November.String(),
	}
	req := &blogpb.UpdateBlogRequest{
		Blog: blog,
	}

	result, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not update blog with error %v", err)
	}

	fmt.Printf("The result was %s \n", result.GetBlog())
}

func (*client) deleteBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.DeleteBlogRequest{
		BlogId: "60a582cf88071f794ab1cd7e",
	}

	result, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not delete blog with error %v", err)
	}

	fmt.Printf("The blog was deleted %v \n", result.GetDeleted())
}

func (*client) listBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.ListBlogRequest{
		AuthorId: "jose",
	}
	result, err := c.ListBlog(context.Background(), req)
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
		fmt.Printf("blog: %s \n", msg.GetBlog())
	}
}

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close() // with defer this line executes at the end of program

	c := blogpb.NewBlogServiceClient(cc)
	//fmt.Printf("create4d client: %f", c)
	client := client{}

	client.createBlog(c)
	client.ReadBlog(c)
	client.UpdateBlog(c)
	client.deleteBlog(c)
	client.listBlog(c)
}
