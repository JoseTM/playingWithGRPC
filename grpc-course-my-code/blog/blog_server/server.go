package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"../blogpb"
)

var collection *mongo.Collection

type server struct{}
type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (*server) CreateBlog(ctx context.Context, in *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := in.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("somthing was wrong inserting in mongo: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintln("can not parse OID "),
		)
	}
	blog.Id = oid.Hex()
	response := &blogpb.CreateBlogResponse{
		Blog: blog,
	}
	return response, nil
}

func (*server) ReadBlog(ctx context.Context, in *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	blogID, err := primitive.ObjectIDFromHex(in.GetBlogId())
	if err != nil {
		log.Fatalf("get an error parsing data from string to object: %v", err)
	}
	filter := bson.M{"_id": blogID}
	res := collection.FindOne(context.Background(), filter)
	if res.Err() != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("somthing was wrong reading from mongo: %v", res.Err()),
		)
	}

	blogItem := &blogItem{}
	errorDecoding := res.Decode(blogItem)
	if errorDecoding != nil {
		fmt.Printf("<< Error Decoding >> : %v\n", errorDecoding)
	}
	blog := &blogpb.Blog{
		Id:       blogItem.ID.Hex(),
		AuthorId: blogItem.AuthorID,
		Title:    blogItem.Title,
		Content:  blogItem.Content,
	}

	response := &blogpb.ReadBlogResponse{
		Blog: blog,
	}
	return response, nil
}

func (*server) UpdateBlog(ctx context.Context, in *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := in.GetBlog()
	blogId, errParsing := primitive.ObjectIDFromHex(blog.GetId())
	if errParsing != nil {
		log.Fatalf("get an error parsing data from string to object: %v", errParsing)
	}
	data := blogItem{
		ID:       blogId,
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	filter := bson.M{"_id": data.ID}
	update := bson.M{"$set": data}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("somthing was wrong updating in mongo: %v", err),
		)
	}

	count := res.ModifiedCount
	if count <= 0 {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintln("nothing has been updated"),
		)
	}

	read := collection.FindOne(context.Background(), filter)
	if read.Err() != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("somthing was wrong reading from mongo: %v", read.Err()),
		)
	}

	blogItem := &blogItem{}
	errorDecoding := read.Decode(blogItem)
	if errorDecoding != nil {
		fmt.Printf("<< Error Decoding >> : %v\n", errorDecoding)
	}
	blogResponse := &blogpb.Blog{
		Id:       blogItem.ID.Hex(),
		AuthorId: blogItem.AuthorID,
		Title:    blogItem.Title,
		Content:  blogItem.Content,
	}
	response := &blogpb.UpdateBlogResponse{
		Blog: blogResponse,
	}

	return response, nil

}

func (*server) DeleteBlog(ctx context.Context, in *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogID, err := primitive.ObjectIDFromHex(in.GetBlogId())
	if err != nil {
		log.Fatalf("get an error parsing data from string to object: %v", err)
	}
	filter := bson.M{"_id": blogID}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("somthing was wrong deleting from mongo: %v", err),
		)
	}

	result := (res.DeletedCount > 0)

	response := &blogpb.DeleteBlogResponse{
		Deleted: result,
	}
	return response, nil
}

func (*server) ListBlog(in *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	authorID := in.GetAuthorId()

	filter := bson.M{"author_id": authorID}
	res, err := collection.Find(context.Background(), filter)
	if err != nil {
		return status.Errorf(
			codes.NotFound,
			fmt.Sprintf("somthing was wrong reading from mongo: %v", err),
		)
	}
	defer res.Close(context.Background())

	blogItem := &blogItem{}
	for res.Next(context.Background()) {
		errorDecoding := res.Decode(blogItem)
		if errorDecoding != nil {
			fmt.Printf("<< Error Decoding >> : %v\n", errorDecoding)
		}
		blog := &blogpb.Blog{
			Id:       blogItem.ID.Hex(),
			AuthorId: blogItem.AuthorID,
			Title:    blogItem.Title,
			Content:  blogItem.Content,
		}

		response := &blogpb.ListBlogResponse{
			Blog: blog,
		}
		stream.Send(response)
	}

	return nil
}

func main() {
	//if we crash the go code we get file name an the line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog service started")

	fmt.Println("starting mongodb connection")
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017/?connect=direct")
	client, errCon := mongo.Connect(context.Background(), clientOpts)
	if errCon != nil {
		log.Fatalf("error connecting mongo: %v", errCon)
	}
	_ = client

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("starting Server ...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve blog: %v", err)
		}
	}()

	//wait for Control+C to exit
	ch := make(chan os.Signal, syscall.SIGHUP)
	signal.Notify(ch, os.Interrupt)

	//Block until a ginal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("closing mongodb connection")
	errorDiscon := client.Disconnect(context.Background())
	if errorDiscon != nil {
		log.Fatalf("error disconnecting mongodb: %v", errorDiscon)
	}
	fmt.Println("End of execution")

}
