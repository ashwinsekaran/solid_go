package main

import (
	"context"
	"fmt"
	"log"
	"solid_go/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewInvoiceServiceClient(conn)
	resp, err := client.CreateInvoice(context.Background(), &proto.Invoice{Customer: "John", Amount: 100.50})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s", resp.Message)
}
