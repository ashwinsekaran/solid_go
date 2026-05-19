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
	// Dial opens a persistent connection to the gRPC server.
	// insecure.NewCredentials() skips TLS — fine for local dev, not for production.
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewInvoiceServiceClient(conn)

	// ── Unary RPC ─────────────────────────────────────────────────────────────
	// One request, one response — the call blocks until the server replies.
	resp, err := client.CreateInvoice(context.Background(), &proto.Invoice{Customer: "John", Amount: 100.50})
	if err != nil {
		log.Fatalf("could not create invoice: %v", err)
	}
	fmt.Printf("Unary response: %s \n", resp.Message)

	// ── Server-streaming RPC ──────────────────────────────────────────────────
	// Client sends one request; server streams back multiple responses.
	// Recv() blocks for each message; io.EOF signals the stream is done.
	resp2, err := client.ListInvoices(context.Background(), &proto.Invoice{Customer: "John"})
	if err != nil {
		log.Fatalf("could not list invoices: %v", err)
	}
	for {
		respx, err := resp2.Recv()
		if err == io.EOF {
			break // server finished streaming
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}
		fmt.Printf("Server stream: %s\n", respx.Message)
	}

	// ── Client-streaming RPC ──────────────────────────────────────────────────
	// Client opens a stream and sends multiple messages; CloseAndRecv() signals
	// end-of-stream and blocks waiting for the server's single aggregate response.
	respy, err := client.UploadInvoices(context.Background())
	if err != nil {
		log.Fatalf("could not open upload stream: %v", err)
	}
	for i := 1; i <= 5; i++ {
		err := respy.Send(&proto.Invoice{Customer: "John", Amount: float64(i)})
		if err != nil {
			log.Fatalf("send error: %v", err)
		}
	}
	respz, err := respy.CloseAndRecv() // tell server we're done and await its summary
	if err != nil {
		log.Fatalf("close and recv error: %v", err)
	}
	fmt.Printf("Client stream response: %s\n", respz.Message)

	// ── Bidirectional-streaming RPC ───────────────────────────────────────────
	// Both sides send and receive concurrently.
	// A goroutine handles sending while the main goroutine reads responses.
	bdrs, err := client.SyncInvoices(context.Background())
	if err != nil {
		log.Fatalf("could not sync: %v", err)
	}

	// Send goroutine: push 5 invoices then signal we're done with CloseSend().
	go func() {
		for i := 1; i <= 5; i++ {
			err := bdrs.Send(&proto.Invoice{Customer: "John", Amount: float64(i)})
			if err != nil {
				log.Fatalf("bidi send error: %v", err)
			}
		}
		bdrs.CloseSend() // half-close: server can still send after we stop
	}()

	// Receive responses until the server closes its side of the stream.
	for {
		resp, err := bdrs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("bidi recv error: %v", err)
		}
		fmt.Printf("Bidi response: %s\n", resp.Message)
	}
}
