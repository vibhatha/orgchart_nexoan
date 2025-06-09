package main

import (
	"fmt"
	"log"
	"os"

	"orgchart_nexoan/api"
)

func main() {
	// Get the data directory from command line argument or use default
	dataDir := "data/gota_gazettes/2020-01-01"
	if len(os.Args) > 1 {
		dataDir = os.Args[1]
	}

	// Create API client
	client := api.NewClient("http://localhost:8080/entities", "http://localhost:8081/v1/entities")

	//Create government node first
	// government, err := client.CreateGovernmentNode()
	// if err != nil {
	// 	log.Fatalf("Failed to create government node: %v", err)
	// }
	// fmt.Printf("Successfully created government node with ID: %s\n", government.ID)

	// err = client.ProcessTransactions(dataDMir)
	// if err != nil {
	// 	log.Fatalf("Failed to process transactions: %v", err)
	// }

	// Process transactions
	err := client.ProcessTransactions(dataDir)
	if err != nil {
		log.Fatalf("Failed to process transactions: %v", err)
	}

	fmt.Println("Successfully processed all transactions")
}
