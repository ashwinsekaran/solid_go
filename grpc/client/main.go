package main

import (
	"context"
	"fmt"
	"io"
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
	fmt.Printf("Greeting: %s \n", resp.Message)

	resp2, err := client.ListInvoices(context.Background(), &proto.Invoice{Customer: "John"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		respx, err := resp2.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Printf("Greeting: %s\n", respx.Message)
	}

	respy, err := client.UploadInvoices(context.Background())

	for i := 1; i <= 5; i++ {
		err := respy.Send(&proto.Invoice{Customer: "John", Amount: float64(i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

	}

	respz, err := respy.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s\n", respz.Message)

	bdrs, err := client.SyncInvoices(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			err := bdrs.Send(&proto.Invoice{Customer: "John", Amount: float64(i)})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			// tell server you're done sendin
		} // g
		bdrs.CloseSend()
	}()

	for {
		resp, err := bdrs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("Greeting: %s\n", resp.Message)
	}

}
