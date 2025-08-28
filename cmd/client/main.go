package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("new client: %v\n", err)
	}

	client := newsv1.NewNewsServiceClient(conn)
	ctx := context.Background()

	// res, err := client.CreateNews(ctx, &newsv1.CreateNewsRequest{
	// 	Id:      uuid.NewString(),
	// 	Author:  "Nabin Dhami",
	// 	Title:   "Random Title",
	// 	Summary: "Random Summary",
	// 	Content: "Random Content",
	// 	Source:  "https://google.com",
	// 	Tags:    []string{"Title", "Description"},
	// })
	// if err != nil {
	// 	log.Fatalf("create news: %v", err)
	// }

	// // log.Print(res)

	// result, err := client.GetNews(ctx, &newsv1.GetNewsRequest{Id: res.Id})
	// if err != nil {
	// 	log.Fatalf("create news: %v", err)
	// }
	// log.Print(result)

	for i := range 2 {
		_, err := client.CreateNews(ctx, &newsv1.CreateNewsRequest{
			Id:      uuid.NewString(),
			Author:  fmt.Sprintf("Nabin Dhami %d", i),
			Title:   "Random Title",
			Summary: "Random Summary",
			Content: "Random Content",
			Source:  "https://google.com",
			Tags:    []string{"Title", "Description"},
		})
		if err != nil {
			log.Fatalf("create news: %v", err)
		}
	}

	getAllRes, err := client.GetAllNews(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("get all news: %v", err)
	}

	allNews := make([]*newsv1.GetNewsResponse, 0)
	for {
		res, err := getAllRes.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("get all news stream: %v", err)
		}
		allNews = append(allNews, res)
		println("MESSAGE RECEIVED!")
	}
	log.Printf("all news: %v", allNews)

	var clientStream grpc.ClientStreamingClient[newsv1.CreateNewsRequest, emptypb.Empty]
	for i, n := range allNews {
		clientStream, err = client.UpdateNews(ctx)
		if err != nil {
			log.Fatalf("update news sream: %v", err)
		}

		clientStream.Send(&newsv1.CreateNewsRequest{
			Id:      n.Id,
			Author:  n.Author + fmt.Sprintf("updated %d", i),
			Title:   n.Title,
			Tags:    n.Tags,
			Summary: n.Summary,
			Content: n.Content,
			Source:  n.Source,
		})
		if err != nil {
			log.Fatalf("update news send: %v", err)
		}
	}

	if _, err := clientStream.CloseAndRecv(); err != nil {
		log.Fatalf("client stream close: %v", err)
	}

	updatedNews := make([]*newsv1.GetNewsResponse, 0)
	for {
		res, err := getAllRes.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("get all news stream: %v", err)
		}
		updatedNews = append(updatedNews, res)
		println("MESSAGE RECEIVED!")
	}
	log.Printf("updated news: %v", updatedNews)
}
