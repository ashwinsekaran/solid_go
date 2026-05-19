package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"solid_go/grpc/proto"

	"google.golang.org/grpc"
)

// InvoiceServer implements all four RPC patterns defined in the proto service.
// Embedding UnimplementedInvoiceServiceServer satisfies the interface for any
// methods we haven't implemented yet, future-proofing against proto additions.
type InvoiceServer struct {
	proto.UnimplementedInvoiceServiceServer
}

// CreateInvoice is a unary RPC: one request in, one response out.
// This is the simplest gRPC pattern — equivalent to a normal function call over the network.
func (i *InvoiceServer) CreateInvoice(ctx context.Context, req *proto.Invoice) (*proto.InvoiceResponse, error) {
	fmt.Printf("Received request for invoice for customer: %s\n", req.Customer)
	return &proto.InvoiceResponse{
		Message: "Invoice created for customer: " + req.Customer,
	}, nil
}

// UploadInvoices is a client-streaming RPC: the client sends a stream of Invoice messages,
// and the server replies with a single response once the stream ends (io.EOF).
func (i *InvoiceServer) UploadInvoices(stream proto.InvoiceService_UploadInvoicesServer) error {
	for {
		resp, err := stream.Recv() // block until the next message arrives
		if err == io.EOF {
			// Client has finished sending — send the single aggregate response.
			return stream.SendAndClose(&proto.InvoiceResponse{Message: "All invoices uploaded"})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received request for invoice for customer: %s\n", resp.Customer)
	}
}

// ListInvoices is a server-streaming RPC: the client sends one request and the server
// streams back multiple responses. The client iterates them with Recv() until io.EOF.
func (i *InvoiceServer) ListInvoices(req *proto.Invoice, stream proto.InvoiceService_ListInvoicesServer) error {
	fmt.Printf("Received request for list of invoices for customer: %s\n", req.Customer)
	for idx := 1; idx <= 5; idx++ {
		err := stream.Send(&proto.InvoiceResponse{Message: fmt.Sprintf("Invoice %d", idx)})
		if err != nil {
			return err // client disconnected or network error
		}
	}
	return nil // returning nil closes the server-side stream
}

// SyncInvoices is a bidirectional-streaming RPC: both sides send and receive concurrently.
// Here the server echoes an acknowledgement for every invoice the client sends.
func (i *InvoiceServer) SyncInvoices(stream proto.InvoiceService_SyncInvoicesServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil // client is done; closing the handler closes the server stream
		}
		if err != nil {
			return err
		}
		println(resp.Customer)
		// Send a response for each received message — true bidirectional exchange.
		err = stream.Send(&proto.InvoiceResponse{
			Message: "Synced invoice for: " + resp.Customer,
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	// Register our implementation so the gRPC server knows which struct handles each RPC.
	proto.RegisterInvoiceServiceServer(server, &InvoiceServer{})

	err = server.Serve(listen) // blocks, accepting connections
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
