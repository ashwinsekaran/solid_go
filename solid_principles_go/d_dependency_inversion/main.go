package main

import "fmt"

type Invoice struct {
	Customer string
	Amount   float64
}

// InvoiceRepository is the abstraction that the high-level service depends on.
// DIP: InvoiceService never imports a concrete DB package — it only knows this interface.
type InvoiceRepository interface {
	Save(inv Invoice)
	GetById(id float64) Invoice
}

// InvoiceService is the high-level module. It holds a repository interface, not a struct.
// Swapping the storage backend requires no change here.
type InvoiceService struct {
	repo InvoiceRepository // injected at construction time
}

// MySQLRepository and MongoDBRepository are low-level modules that satisfy the interface.
type MySQLRepository struct{}
type MongoDBRepository struct{}

func (m *MySQLRepository) Save(inv Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func (m *MongoDBRepository) Save(inv Invoice) {
	fmt.Printf("Saving invoice for %s: %f\n", inv.Customer, inv.Amount)
}

func (m *MySQLRepository) GetById(id float64) Invoice {
	return Invoice{Customer: "John", Amount: 100.50}
}

func (m *MongoDBRepository) GetById(id float64) Invoice {
	return Invoice{Customer: "John", Amount: 100.50}
}

// Create delegates to whatever repository was injected — MySQL or Mongo, it doesn't matter.
func (s *InvoiceService) Create(inv Invoice) {
	s.repo.Save(inv)
}

func (s *InvoiceService) GetById(id float64) Invoice {
	return s.repo.GetById(id)
}

func main() {
	invoice := Invoice{Customer: "John", Amount: 100.50}

	// Inject MySQLRepository — InvoiceService has no idea it's MySQL.
	mysqlService := InvoiceService{repo: &MySQLRepository{}}
	mysqlService.Create(invoice)
	data := mysqlService.GetById(1)
	fmt.Printf("Got sql invoice: %s\n", data.Customer)

	// Swap the dependency for MongoDB — InvoiceService code is unchanged.
	mongoService := InvoiceService{repo: &MongoDBRepository{}}
	mongoService.Create(invoice)
	inv := mongoService.GetById(1)
	fmt.Printf("Got mongo invoice: %s\n", inv.Customer)
}
