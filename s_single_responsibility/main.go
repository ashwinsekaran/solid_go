package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	PrintInvoice(invoice)
	SaveInvoice(invoice)
	SendEmail(invoice)
}

func PrintInvoice(invoice Invoice) {
	fmt.Printf("Invoice for %s: %f\n", invoice.Customer, invoice.Amount)
}

func SaveInvoice(invoice Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", invoice.Customer, invoice.Amount)
}

func SendEmail(invoice Invoice) {
	fmt.Printf("Sending email to %s for invoice: %f\n", invoice.Customer, invoice.Amount)
}
