package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("new client: %v\n", err)
	}

	client := newsv1.NewNewsServiceClient(conn)
	ctx := context.Background()

	res, err := client.CreateNews(ctx, &newsv1.CreateNewsRequest{
		Id:      uuid.NewString(),
		Author:  "Nabin Dhami",
		Title:   "Random Title",
		Summary: "Random Summary",
		Content: "Random Content",
		Source:  "https://google.com",
		Tags:    []string{"Title", "Description"},
	})
	if err != nil {
		log.Fatalf("create news: %v", err)
	}

	// log.Print(res)

	result, err := client.GetNews(ctx, &newsv1.GetNewsRequest{Id: res.Id})
	if err != nil {
		log.Fatalf("create news: %v", err)
	}
	log.Print(result)
}
