package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

type InvoiceRepository interface {
	Save(inv Invoice)
	GetById(id float64) Invoice
}

type InvoiceService struct {
	repo InvoiceRepository
}

type MySQLRepository struct {
}

type MongoDBRepository struct {
}

func (m *MySQLRepository) Save(inv Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func (m *MongoDBRepository) Save(inv Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func (m *MySQLRepository) GetById(id float64) Invoice {
	return Invoice{
		Customer: "John",
		Amount:   100.50,
	}
}

func (m *MongoDBRepository) GetById(id float64) Invoice {
	return Invoice{
		Customer: "John",
		Amount:   100.50,
	}
}

func (s *InvoiceService) Create(inv Invoice) {
	s.repo.Save(inv)
}

func (s *InvoiceService) GetById(id float64) Invoice {
	return s.repo.GetById(id)
}

func main() {
	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
	}

	mysqlService := InvoiceService{
		repo: &MySQLRepository{},
	}
	mysqlService.Create(invoice)
	data := mysqlService.GetById(1)
	fmt.Printf("Got sql invoice: %s\n", data.Customer)

	mongoService := InvoiceService{
		repo: &MongoDBRepository{},
	}
	mongoService.Create(invoice)
	inv := mongoService.GetById(1)
	fmt.Printf("Got mongo invoice: %s\n", inv.Customer)

}
