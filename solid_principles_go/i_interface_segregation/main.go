package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

// InvoicePrinter is a narrow interface: only types that print need to implement it.
// ISP: splitting one fat interface into focused ones means no type carries unused methods.
type InvoicePrinter interface {
	Print(inv Invoice)
}

// InvoiceExporter is a separate narrow interface for types that export invoices.
type InvoiceExporter interface {
	Export(inv Invoice)
}

// BasicReporter can only print — it implements InvoicePrinter but NOT InvoiceExporter.
// ISP: it isn't forced to have an Export method it doesn't need.
type BasicReporter struct{}

// FileExporter can only export — it implements InvoiceExporter but NOT InvoicePrinter.
type FileExporter struct{}

// FullExporter implements both interfaces — it can print and export.
type FullExporter struct{}

// HandlePrint accepts anything that can print — it doesn't care about exporting.
func HandlePrint(inv Invoice, printer InvoicePrinter) {
	printer.Print(inv)
}

// HandleExport accepts anything that can export — it doesn't care about printing.
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

	// FullExporter satisfies both interfaces — pass it to either handler.
	HandleExport(invoice, &FullExporter{})
	HandlePrint(invoice, &FullExporter{})

	// BasicReporter only satisfies InvoicePrinter — passing it to HandleExport would not compile.
	HandlePrint(invoice, &BasicReporter{})

	// FileExporter only satisfies InvoiceExporter — passing it to HandlePrint would not compile.
	HandleExport(invoice, &FileExporter{})
}
