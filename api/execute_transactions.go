package api

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ProcessTransactions processes all transactions from CSV files in the specified directory
func (c *Client) ProcessTransactions(dataDir string, processType string) error {
	// Initialize entity counters based on process type
	var entityCounters map[string]int
	if processType == "organisation" {
		entityCounters = map[string]int{
			"minister":   0,
			"department": 0,
		}
	} else if processType == "person" {
		entityCounters = map[string]int{
			"citizen": 0,
		}
	} else {
		return fmt.Errorf("invalid process type: %s", processType)
	}

	// Get all CSV files in the directory
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dataDir, err)
	}

	// Collect all transactions from all files
	var allTransactions []map[string]interface{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			// Check if file contains ADD in its name (either as _ADD.csv or ADD.csv)
			fileName := strings.TrimSuffix(file.Name(), ".csv")
			if fileName == "ADD" || strings.HasSuffix(fileName, "_ADD") {
				// Extract file type from filename (e.g., "ADD" from "2403-38_ADD.csv" or "ADD.csv")
				fileType := "ADD" // All files with ADD are considered ADD transactions

				// Load transactions from the CSV file
				transactions, err := loadTransactions(filepath.Join(dataDir, file.Name()), fileType)
				if err != nil {
					return fmt.Errorf("failed to load transactions from %s: %w", file.Name(), err)
				}
				allTransactions = append(allTransactions, transactions...)
			}
		}
	}

	fmt.Printf("Number of transactions: %d\n", len(allTransactions))

	// Sort transactions by transaction_id, handling numeric parts correctly
	sort.Slice(allTransactions, func(i, j int) bool {
		idI := allTransactions[i]["transaction_id"].(string)
		idJ := allTransactions[j]["transaction_id"].(string)

		// Split the IDs into parts
		partsI := strings.Split(idI, "_")
		partsJ := strings.Split(idJ, "_")

		// Compare the first part (e.g., "2153/12")
		if partsI[0] != partsJ[0] {
			return partsI[0] < partsJ[0]
		}

		// Compare the second part (e.g., "tr")
		if partsI[1] != partsJ[1] {
			return partsI[1] < partsJ[1]
		}

		// Compare the numeric part by converting to integers
		numI := strings.TrimPrefix(partsI[2], "tr_")
		numJ := strings.TrimPrefix(partsJ[2], "tr_")

		// Convert to integers for numeric comparison
		valI, _ := strconv.Atoi(numI)
		valJ, _ := strconv.Atoi(numJ)
		return valI < valJ
	})

	// Process transactions in order
	for _, transaction := range allTransactions {
		// For now, only process ADD transactions
		if transaction["file_type"] == "ADD" {
			// Check if the transaction type matches the process type
			childType := transaction["child_type"].(string)
			fmt.Printf("Processing transaction: %s\n", transaction["transaction_id"])
			fmt.Printf("Child type: %s\n", childType)
			fmt.Printf("Process type: %s\n", processType)

			if (processType == "organisation" && (childType == "minister" || childType == "department")) ||
				(processType == "person" && childType == "citizen") {
				var newCounter int
				var err error

				if processType == "person" && childType == "citizen" {
					newCounter, err = c.AddPersonEntity(transaction, entityCounters)
				} else {
					newCounter, err = c.AddOrgEntity(transaction, entityCounters)
				}

				if err != nil {
					return fmt.Errorf("failed to process add transaction %s: %w", transaction["transaction_id"], err)
				}
				entityCounters[childType] = newCounter
				fmt.Printf("Processed Add transaction: %s\n", transaction["transaction_id"])
			} else {
				fmt.Printf("Skipping transaction %s: type %s does not match process type %s\n",
					transaction["transaction_id"], childType, processType)
			}
		} else {
			fmt.Printf("Skipping file: %s\n", transaction["transaction_id"])
		}
	}

	return nil
}

// loadTransactions reads and processes transactions from a CSV file
func loadTransactions(filePath string, fileType string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header from %s: %w", filePath, err)
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from %s: %w", filePath, err)
	}

	var transactions []map[string]interface{}
	// Process each record
	for _, record := range records {
		transaction := make(map[string]interface{})
		for i, value := range record {
			transaction[header[i]] = value
		}
		// Add file type to transaction
		transaction["file_type"] = fileType
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
