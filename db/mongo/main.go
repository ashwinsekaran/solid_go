package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Invoice struct {
	Customer string
	Amount   float64
	Status   string
}

func main() {
	// Use a timeout context for the connection — prevents hanging forever if Mongo is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Ping verifies the driver can actually reach the server (Connect is lazy).
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	println("Connected to MongoDB")

	invoice := Invoice{Customer: "John", Amount: 100.50, Status: "Paid"}

	// Collection is the Mongo equivalent of an SQL table.
	collection := client.Database("invoicedb").Collection("invoices")

	// Drop for a clean slate on each run — useful for demos, not for production.
	collection.Drop(ctx)

	// InsertOne returns a result containing the auto-generated _id.
	res, err := collection.InsertOne(ctx, invoice)
	fmt.Println(res.InsertedID)

	// ── Read ──────────────────────────────────────────────────────────────────
	// bson.M{"customer": "John"} is the filter — equivalent to WHERE customer = 'John'.
	var result Invoice
	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// ── FindOneAndUpdate (atomic read-modify) ─────────────────────────────────
	// $set updates only the specified fields — other fields are left unchanged.
	collection.FindOneAndUpdate(ctx, bson.M{"customer": "John"}, bson.M{"$set": bson.M{"status": "Pending"}})

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// ── UpdateOne ─────────────────────────────────────────────────────────────
	_, err = collection.UpdateOne(ctx, bson.M{"customer": "John"}, bson.M{"$set": bson.M{"amount": 200.0}})
	if err != nil {
		return
	}

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// ── DeleteOne ─────────────────────────────────────────────────────────────
	_, err = collection.DeleteOne(ctx, bson.M{"customer": "John"})
	if err != nil {
		return
	}

	// ErrNoDocuments is the sentinel error for a FindOne that matches nothing.
	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("Document not found")
	} else if err != nil {
		panic(err)
	}

	// ── InsertMany ────────────────────────────────────────────────────────────
	_, err = collection.InsertMany(ctx, []interface{}{
		Invoice{Customer: "John", Amount: 100.50, Status: "Paid"},
		Invoice{Customer: "John", Amount: 200.0, Status: "Pending"},
		Invoice{Customer: "Jane", Amount: 150.0, Status: "Paid"},
		Invoice{Customer: "Bob", Amount: 300.0, Status: "Pending"},
	})
	if err != nil {
		return
	}

	// ── Compound index ────────────────────────────────────────────────────────
	// Value: 1 = ascending. The order of fields matters for index usage:
	// a query on {customer, status} can use this index; a query on {status} alone cannot
	// (it doesn't start with the leading field).
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "customer", Value: 1},
			{Key: "status", Value: 1},
		},
	})
	if err != nil {
		return
	}

	// This query uses the compound index (customer is the leading field).
	cursor, err := collection.Find(ctx, bson.M{"customer": "John", "status": "Paid"})
	var results []Invoice
	err = cursor.All(ctx, &results)
	if err != nil {
		return
	}
	fmt.Printf("First result - %v\n", results)

	// This query filters only on status — the leading field is missing,
	// so MongoDB cannot use the compound index and falls back to a collection scan.
	cursor1, err := collection.Find(ctx, bson.M{"status": "Paid"})
	var results1 []Invoice
	err = cursor1.All(ctx, &results1)
	if err != nil {
		return
	}
	fmt.Printf("Second result - %v\n", results1)

	// Explain plans show exactly which execution strategy MongoDB chose.
	fmt.Println("\n--- Explain Plan: Query WITH index (customer + status) ---")
	runExplain(ctx, client, bson.M{"customer": "John", "status": "Paid"})

	fmt.Println("\n--- Explain Plan: Query WITHOUT leading index field (status only) ---")
	runExplain(ctx, client, bson.M{"status": "Paid"})
}

// runExplain executes an explain command and prints a concise summary of the winning plan.
// Key fields to watch: stage (IXSCAN = index used, COLLSCAN = full table scan),
// totalDocsExamined (lower is better), and executionTimeMillis.
func runExplain(ctx context.Context, client *mongo.Client, filter bson.M) {
	command := bson.D{
		{Key: "explain", Value: bson.D{
			{Key: "find", Value: "invoices"},
			{Key: "filter", Value: filter},
		}},
		{Key: "verbosity", Value: "executionStats"},
	}

	var explainResult bson.M
	err := client.Database("invoicedb").RunCommand(ctx, command).Decode(&explainResult)
	if err != nil {
		panic(err)
	}

	execStats := explainResult["executionStats"].(bson.M)
	winningPlan := explainResult["queryPlanner"].(bson.M)["winningPlan"]

	summary := bson.M{
		"stage":               winningPlan,
		"totalDocsExamined":   execStats["totalDocsExamined"],
		"totalKeysExamined":   execStats["totalKeysExamined"],
		"executionTimeMillis": execStats["executionTimeMillis"],
		"nReturned":           execStats["nReturned"],
	}

	out, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(out))
}
