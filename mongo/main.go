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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	println("Connected to MongoDB")

	invoice := Invoice{
		Customer: "John",
		Amount:   100.50,
		Status:   "Paid",
	}

	collection := client.Database("invoicedb").Collection("invoices")

	err = collection.Drop(ctx)
	if err != nil {
		return
	}

	res, err := collection.InsertOne(ctx, invoice)

	fmt.Println(res.InsertedID)

	var result Invoice

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	collection.FindOneAndUpdate(ctx, bson.M{"customer": "John"}, bson.M{"$set": bson.M{"status": "Pending"}})

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	_, err = collection.UpdateOne(ctx, bson.M{"customer": "John"}, bson.M{"$set": bson.M{"amount": 200.0}})
	if err != nil {
		return
	}

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	_, err = collection.DeleteOne(ctx, bson.M{"customer": "John"})
	if err != nil {
		return
	}

	err = collection.FindOne(ctx, bson.M{"customer": "John"}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("Document not found")
	} else if err != nil {
		panic(err)
	}

	_, err = collection.InsertMany(ctx, []interface{}{
		Invoice{Customer: "John", Amount: 100.50, Status: "Paid"},
		Invoice{Customer: "John", Amount: 200.0, Status: "Pending"},
		Invoice{Customer: "Jane", Amount: 150.0, Status: "Paid"},
		Invoice{Customer: "Bob", Amount: 300.0, Status: "Pending"},
	})
	if err != nil {
		return
	}
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "customer", Value: 1},
			{Key: "status", Value: 1},
		},
	})
	if err != nil {
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"customer": "John", "status": "Paid"})
	var results []Invoice
	err = cursor.All(ctx, &results)
	if err != nil {
		return
	}

	fmt.Printf("First result - %v\n", results)

	cursor1, err := collection.Find(ctx, bson.M{"status": "Paid"})
	var results1 []Invoice
	err = cursor1.All(ctx, &results1)
	if err != nil {
		return
	}

	fmt.Printf("Second result - %v\n", results1)

	// Explain plan — Query 1: uses compound index (customer + status)
	fmt.Println("\n--- Explain Plan: Query WITH index (customer + status) ---")
	runExplain(ctx, client, bson.M{"customer": "John", "status": "Paid"})

	// Explain plan — Query 2: missing leading field, cannot use compound index
	fmt.Println("\n--- Explain Plan: Query WITHOUT leading index field (status only) ---")
	runExplain(ctx, client, bson.M{"status": "Paid"})
}

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

	// Extract key stats only — full output is very verbose
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
