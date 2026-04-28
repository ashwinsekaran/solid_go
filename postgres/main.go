package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Address struct {
	City string `json:"city"`
}

type InvoiceMetadata struct {
	Status  string   `json:"status"`
	Amount  float64  `json:"amount"`
	Tags    []string `json:"tags"`
	Address Address  `json:"address"`
}

type Invoice struct {
	ID       int
	Customer string
	Metadata InvoiceMetadata
}

func main() {
	connStr := "postgres://admin:password@localhost:5432/invoicedb?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS invoices (
					id serial PRIMARY KEY,
					customer VARCHAR(255) NOT NULL,
					metadata JSONB NOT NULL)`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Table created!")

	invoices := []Invoice{
		{
			Customer: "John",
			Metadata: InvoiceMetadata{
				Status:  "Paid",
				Amount:  100.50,
				Tags:    []string{"urgent", "international"},
				Address: Address{City: "Berlin"},
			},
		},
		{
			Customer: "Jane",
			Metadata: InvoiceMetadata{
				Status:  "Pending",
				Amount:  200.0,
				Tags:    []string{"domestic"},
				Address: Address{City: "Munich"},
			},
		},
		{
			Customer: "Bob",
			Metadata: InvoiceMetadata{
				Status:  "Paid",
				Amount:  300.0,
				Tags:    []string{"urgent"},
				Address: Address{City: "Berlin"},
			},
		},
	}

	for _, meta := range invoices {
		metaJson, err := json.Marshal(meta.Metadata)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("INSERT INTO invoices (customer, metadata) VALUES ($1, $2)", meta.Customer, metaJson)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted invoices!")

	status := "Paid"
	amount := 100.50
	city := "Berlin"
	rows, err := db.Query((`
SELECT id, customer, metadata FROM invoices 
where metadata ->> 'status' = $1 
AND (metadata ->> 'amount')::float = $2 
AND metadata #>> '{address, city}' = $3`),
		status, amount, city)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Invoice
	for rows.Next() {
		var invoice Invoice
		var metaJson []byte
		err = rows.Scan(&invoice.ID, &invoice.Customer, &metaJson)
		if err != nil {
			panic(err)
		}
		fmt.Printf("invoice id: %d, customer: %s\n", invoice.ID, invoice.Customer)
		err = json.Unmarshal(metaJson, &invoice.Metadata)
		if err != nil {
			panic(err)
		}
		result = append(result, invoice)

		for _, inv := range result {
			fmt.Printf("ID: %d, Customer: %s, Status: %s, Amount: %.2f, City: %s\n",
				inv.ID, inv.Customer, inv.Metadata.Status, inv.Metadata.Amount, inv.Metadata.Address.City)
		}
	}
}
