package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

// DiscountStrategy is the extension point.
// OCP: to add a new discount type, implement this interface — don't modify PrintFinalAmount.
type DiscountStrategy interface {
	Apply(amount float64) float64
}

// PercentageDiscount reduces the amount by a percentage of itself.
type PercentageDiscount struct {
	Percentage float64
}

// FlatDiscount subtracts a fixed amount regardless of the invoice total.
type FlatDiscount struct {
	Amount float64
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	// Pass different strategies without changing PrintFinalAmount — open for extension.
	discount := PercentageDiscount{Percentage: 50}
	PrintFinalAmount(invoice, discount)

	discount2 := FlatDiscount{Amount: 10}
	PrintFinalAmount(invoice, discount2)
}

// Apply reduces amount by the given percentage.
func (p PercentageDiscount) Apply(amount float64) float64 {
	return amount - (amount * p.Percentage / 100)
}

// Apply subtracts the flat discount amount.
func (f FlatDiscount) Apply(amount float64) float64 {
	return amount - f.Amount
}

// PrintFinalAmount is closed for modification — it works with any DiscountStrategy,
// including ones that don't exist yet, without any changes to this function.
func PrintFinalAmount(inv Invoice, discount DiscountStrategy) {
	a := discount.Apply(inv.Amount)
	fmt.Printf("Invoice for %s: %f discounted on total %f \n", inv.Customer, a, inv.Amount)
}
