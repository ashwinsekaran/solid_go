package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

type InvoicePrinter interface {
	Print(inv Invoice)
}

type InvoiceExporter interface {
	Export(inv Invoice)
}

type BasicReporter struct {
}
type FileExporter struct {
}
type FullExporter struct {
}

func HandlePrint(inv Invoice, printer InvoicePrinter) {
	printer.Print(inv)
}

func HandleExport(inv Invoice, exporter InvoiceExporter) {
	exporter.Export(inv)
}

func (b *BasicReporter) Print(inv Invoice) {
	fmt.Printf("Print Invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func (f *FileExporter) Export(inv Invoice) {
	fmt.Printf("Export Invoice for %s to file: %f\n", inv.Customer, inv.Amount)
}

func (f *FullExporter) Export(inv Invoice) {
	fmt.Printf("Full Export Invoice for %s to file: %f\n", inv.Customer, inv.Amount)
}

func (f *FullExporter) Print(inv Invoice) {
	fmt.Printf("Full Print Invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	HandleExport(invoice, &FullExporter{})
	HandlePrint(invoice, &FullExporter{})

	HandlePrint(invoice, &BasicReporter{})

	HandleExport(invoice, &FileExporter{})

	//HandleExport(invoice, &BasicReporter{})

}
