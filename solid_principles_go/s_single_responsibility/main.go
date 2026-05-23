package main

import "fmt"

// Invoice is a plain data carrier — it holds no behaviour.
// SRP: the struct's only job is to represent invoice data.
type Invoice struct {
	Customer string
	Amount   float64
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	// Each call below delegates to a function with a single, distinct responsibility.
	PrintInvoice(invoice)
	SaveInvoice(invoice)
	SendEmail(invoice)
}

// PrintInvoice is responsible for one thing: rendering the invoice to stdout.
// If the output format changes, only this function needs to change.
func PrintInvoice(invoice Invoice) {
	fmt.Printf("Invoice for %s: %f\n", invoice.Customer, invoice.Amount)
}

// SaveInvoice is responsible for one thing: persisting the invoice.
// Changing the storage backend (e.g. DB vs file) only affects this function.
func SaveInvoice(invoice Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", invoice.Customer, invoice.Amount)
}

// SendEmail is responsible for one thing: notifying the customer.
// Email template changes are isolated here and don't ripple to other functions.
func SendEmail(invoice Invoice) {
	fmt.Printf("Sending email to %s for invoice: %f\n", invoice.Customer, invoice.Amount)
}
