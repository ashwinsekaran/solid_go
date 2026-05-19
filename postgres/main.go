package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq" // side-effect import: registers the "postgres" driver with database/sql
)

// Address and InvoiceMetadata are nested structs that will be serialised to/from JSONB.
type Address struct {
	City string `json:"city"`
}

type InvoiceMetadata struct {
	Status  string   `json:"status"`
	Amount  float64  `json:"amount"`
	Tags    []string `json:"tags"`
	Address Address  `json:"address"`
}

// Invoice maps to the invoices table. Metadata is stored as JSONB in PostgreSQL
// and marshalled/unmarshalled in Go using encoding/json.
type Invoice struct {
	ID       int
	Customer string
	Metadata InvoiceMetadata
}

func main() {
	connStr := "postgres://admin:password@localhost:5432/invoicedb?sslmode=disable"

	// sql.Open only validates the driver name and DSN format — it does NOT connect yet.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Ping actually opens a connection and checks that the DB is reachable.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL!")

	// DDL: create the table only if it doesn't already exist — idempotent migration.
	// JSONB stores JSON as a binary blob that PostgreSQL can index and query efficiently.
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
		{Customer: "John", Metadata: InvoiceMetadata{Status: "Paid", Amount: 100.50, Tags: []string{"urgent", "international"}, Address: Address{City: "Berlin"}}},
		{Customer: "Jane", Metadata: InvoiceMetadata{Status: "Pending", Amount: 200.0, Tags: []string{"domestic"}, Address: Address{City: "Munich"}}},
		{Customer: "Bob", Metadata: InvoiceMetadata{Status: "Paid", Amount: 300.0, Tags: []string{"urgent"}, Address: Address{City: "Berlin"}}},
	}

	for _, meta := range invoices {
		// Marshal the nested struct to JSON bytes before inserting into the JSONB column.
		metaJson, err := json.Marshal(meta.Metadata)
		if err != nil {
			panic(err)
		}
		// $1, $2 are positional placeholders — lib/pq uses PostgreSQL-style params, not ?.
		_, err = db.Exec("INSERT INTO invoices (customer, metadata) VALUES ($1, $2)", meta.Customer, metaJson)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted invoices!")

	status := "Paid"
	amount := 100.50
	city := "Berlin"

	// JSONB path operators:
	//   ->>  extracts a top-level field as text
	//   #>>  extracts a nested field via a path array (e.g. '{address, city}')
	// ::float casts the extracted text to a numeric type for comparison.
	rows, err := db.Query(`
		SELECT id, customer, metadata FROM invoices
		WHERE metadata ->> 'status' = $1
		AND (metadata ->> 'amount')::float = $2
		AND metadata #>> '{address, city}' = $3`,
		status, amount, city)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Invoice
	for rows.Next() {
		var invoice Invoice
		var metaJson []byte
		// Scan the JSONB column into a raw []byte, then unmarshal it back into the struct.
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
