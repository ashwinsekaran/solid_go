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

type InvoiceServer struct {
	proto.UnimplementedInvoiceServiceServer
}

func (i *InvoiceServer) CreateInvoice(ctx context.Context, req *proto.Invoice) (*proto.InvoiceResponse, error) {
	fmt.Printf("Received request for invoice for customer: %s\n", req.Customer)
	return &proto.InvoiceResponse{
		Message: "Invoice created for customer: " + req.Customer + "",
	}, nil
}

func (i *InvoiceServer) UploadInvoices(stream proto.InvoiceService_UploadInvoicesServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.InvoiceResponse{Message: "All invoices uploaded"})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received request for invoice for customer: %s\n", resp.Customer)
	}
}

func (i *InvoiceServer) ListInvoices(req *proto.Invoice, stream proto.InvoiceService_ListInvoicesServer) error {
	fmt.Printf("Received request for list of invoices for customer: %s\n", req.Customer)
	for i := 1; i <= 5; i++ {
		err := stream.Send(&proto.InvoiceResponse{Message: fmt.Sprintf("Invoice %d", i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *InvoiceServer) SyncInvoices(stream proto.InvoiceService_SyncInvoicesServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		println(resp.Customer)
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

	proto.RegisterInvoiceServiceServer(server, &InvoiceServer{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
