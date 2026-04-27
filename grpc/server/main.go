package main

import (
	"context"
	"fmt"
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
