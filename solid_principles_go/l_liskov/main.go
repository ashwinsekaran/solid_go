package main

import "fmt"

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	// LSP: PDFExporter and CSVExporter are interchangeable here.
	// ProcessExport doesn't know which concrete type it receives — and doesn't need to.
	ProcessExport(invoice, &PDFExporter{})
	ProcessExport(invoice, &CSVExporter{})
}

type Invoice struct {
	Customer string
	Amount   float64
}

// InvoiceExporter defines the contract every exporter must honour.
// LSP guarantees that any implementor can substitute another without breaking callers.
type InvoiceExporter interface {
	Export(inv Invoice)
}

type PDFExporter struct{}
type CSVExporter struct{}

// ProcessExport works purely against the InvoiceExporter interface.
// Swapping the exporter never requires changing this function.
func ProcessExport(inv Invoice, exporter InvoiceExporter) {
	exporter.Export(inv)
}

// PDFExporter satisfies InvoiceExporter — callers can substitute it for any other exporter.
func (e *PDFExporter) Export(inv Invoice) {
	fmt.Printf("Exporting invoice for %s to PDF\n", inv.Customer)
}

// CSVExporter satisfies InvoiceExporter — callers can substitute it for any other exporter.
func (e *CSVExporter) Export(inv Invoice) {
	fmt.Printf("Exporting invoice for %s to CSV\n", inv.Customer)
}
