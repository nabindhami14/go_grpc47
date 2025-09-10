package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"buf.build/go/protovalidate"
)

const serviceConfig = `{
	"loadBalancingConfig": [{ "round_robin": {} }],
	"methodConfig": [{
		"name": [{
			"method": "Get",
			"service": "news.v1.NewsService"
		}],
		"retryPolicy": {
			"backoffMultiplier": 1.5,
			"initialBackoff": "0.1s",
			"maxAttempts": 5,
			"maxBackoff": "0.5s",
			"retryableStatusCodes": ["INTERNAL","UNAVAILABLE"]
		},
		"timeout": "2s",
		"waitForReady": true
	}]
}`

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				log.Println("client side interceptors")
				return invoker(ctx, method, req, reply, cc, opts...)
			}),
		grpc.WithChainStreamInterceptor(
			func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return streamer(ctx, desc, cc, method, opts...)
			}),
		grpc.WithDefaultServiceConfig(serviceConfig))

	if err != nil {
		log.Fatalf("new client: %v\n", err)
	}

	client := newsv1.NewNewsServiceClient(conn)
	ctx := context.Background()

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("validator initialization :%v", err)
	}

	msg := &newsv1.CreateNewsRequest{
		Id:      uuid.NewString(),
		Author:  "Nabin Dhami",
		Title:   "Random Title",
		Summary: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Earum possimus nostrum minus quos accusamus at laboriosam ut error iste itaque excepturi, quas tempora incidunt eaque voluptatem repudiandae suscipit aliquid optio inventore molestias? Blanditiis veritatis, facilis obcaecati sequi, quibusdam non voluptate eveniet perspiciatis voluptatem nulla, quia nihil. Delectus placeat quaerat molestiae?",
		Content: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Voluptas libero debitis optio temporibus? Voluptatibus suscipit sit at nemo. Nisi quaerat sunt explicabo porro eveniet veritatis magnam laborum, eius nobis iste accusantium dolores reiciendis quisquam? Laborum est sequi consectetur, dolore cumque dolorum, repudiandae a molestias itaque, sint ratione animi autem odit perferendis. Nisi tempora itaque, architecto enim laboriosam et id, nihil culpa sit, debitis voluptates. Facilis eveniet hic voluptatem sunt, dolorem quod odit, natus recusandae obcaecati, dolor cumque? Dolor quos aliquam asperiores repudiandae mollitia deserunt, repellat rem recusandae, incidunt non nobis! Incidunt earum magni labore, odit eligendi iusto? Fugit, voluptate aliquid!",
		Source:  "https://google.com",
		Tags:    []string{"Title", "Description"},
	}
	if err := validator.Validate(msg); err != nil {
		log.Fatalf("validation error: %v", err)
	}

	res, err := client.CreateNews(ctx, msg)
	if err != nil {
		log.Fatalf("create news: %v", err)
	}
	log.Print(res)

	// result, err := client.GetNews(ctx, &newsv1.GetNewsRequest{Id: res.Id})
	// if err != nil {
	// 	log.Fatalf("create news: %v", err)
	// }
	// log.Print(result)

	// for i := range 2 {
	// 	_, err := client.CreateNews(ctx, &newsv1.CreateNewsRequest{
	// 		Id:      uuid.NewString(),
	// 		Author:  fmt.Sprintf("Nabin Dhami %d", i),
	// 		Title:   "Random Title",
	// 		Summary: "Random Summary",
	// 		Content: "Random Content",
	// 		Source:  "https://google.com",
	// 		Tags:    []string{"Title", "Description"},
	// 	})
	// 	if err != nil {
	// 		log.Fatalf("create news: %v", err)
	// 	}
	// }

	// getAllRes, err := client.GetAllNews(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Fatalf("get all news: %v", err)
	// }

	// allNews := make([]*newsv1.GetNewsResponse, 0)
	// for {
	// 	res, err := getAllRes.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	if err != nil {
	// 		log.Fatalf("get all news stream: %v", err)
	// 	}
	// 	allNews = append(allNews, res)
	// 	println("MESSAGE RECEIVED!")
	// }
	// log.Printf("all news: %v", allNews)

	// var clientStream grpc.ClientStreamingClient[newsv1.CreateNewsRequest, emptypb.Empty]
	// for i, n := range allNews {
	// 	clientStream, err = client.UpdateNews(ctx)
	// 	if err != nil {
	// 		log.Fatalf("update news sream: %v", err)
	// 	}

	// 	clientStream.Send(&newsv1.CreateNewsRequest{
	// 		Id:      n.Id,
	// 		Author:  n.Author + fmt.Sprintf("updated %d", i),
	// 		Title:   n.Title,
	// 		Tags:    n.Tags,
	// 		Summary: n.Summary,
	// 		Content: n.Content,
	// 		Source:  n.Source,
	// 	})
	// 	if err != nil {
	// 		log.Fatalf("update news send: %v", err)
	// 	}
	// }

	// if _, err := clientStream.CloseAndRecv(); err != nil {
	// 	log.Fatalf("client stream close: %v", err)
	// }

	// updatedNews := make([]*newsv1.GetNewsResponse, 0)
	// for {
	// 	res, err := getAllRes.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	if err != nil {
	// 		log.Fatalf("get all news stream: %v", err)
	// 	}
	// 	updatedNews = append(updatedNews, res)
	// 	println("MESSAGE RECEIVED!")
	// }
	// log.Printf("updated news: %v", updatedNews)

	// // BIDIRECTIONAL STREAM
	// deleteStream, err := client.DeleteNews(ctx)
	// if err != nil {
	// 	log.Fatalf("delete news stream: %v", err)
	// }

	// waitC := make(chan struct{})
	// go func() {
	// 	for _, news := range allNews {
	// 		err := deleteStream.Send(&newsv1.GetNewsRequest{Id: news.Id})
	// 		if err != nil {
	// 			log.Fatalf("deleting news: %v", err)
	// 		}
	// 	}
	// 	deleteStream.CloseSend()
	// 	close(waitC)
	// }()

	// for {
	// 	_, err := deleteStream.Recv()
	// 	if errors.Is(err, io.EOF) {
	// 		log.Printf("delete stream ended: %v", err)
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("delete stream: %v", err)
	// 	}
	// 	log.Printf("news deleted!")
	// }

	// <-waitC
}
