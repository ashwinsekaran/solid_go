package main

import "fmt"

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	ProcessExport(invoice, &PDFExporter{})
	ProcessExport(invoice, &CSVExporter{})
}

type Invoice struct {
	Customer string
	Amount   float64
}

type InvoiceExporter interface {
	Export(inv Invoice)
}

type PDFExporter struct{}
type CSVExporter struct{}

func ProcessExport(inv Invoice, exporter InvoiceExporter) {
	exporter.Export(inv)
}

func (e *PDFExporter) Export(inv Invoice) {
	fmt.Printf("Exporting invoice for %s to PDF\n", inv.Customer)
}

func (e *CSVExporter) Export(inv Invoice) {
	fmt.Printf("Exporting invoice for %s to CSV\n", inv.Customer)
}
