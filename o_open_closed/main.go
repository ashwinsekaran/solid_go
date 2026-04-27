package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

type DiscountStrategy interface {
	Apply(amount float64) float64
}

type PercentageDiscount struct {
	Percentage float64
}

type FlatDiscount struct {
	Amount float64
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	discount := PercentageDiscount{
		Percentage: 50,
	}
	PrintFinalAmount(invoice, discount)

	discount2 := FlatDiscount{
		Amount: 10,
	}
	PrintFinalAmount(invoice, discount2)

}

func (p PercentageDiscount) Apply(amount float64) float64 {
	f := amount - (amount * p.Percentage / 100)
	return f
}

func (f FlatDiscount) Apply(amount float64) float64 {
	return amount - f.Amount
}

func PrintFinalAmount(inv Invoice, discount DiscountStrategy) {
	a := discount.Apply(inv.Amount)
	fmt.Printf("Invoice for %s: %f discounted on total %f \n", inv.Customer, a, inv.Amount)
}
