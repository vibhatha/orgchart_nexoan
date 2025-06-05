package main

import (
	"fmt"
	"log"
	"time"

	"orgchart_nexoan/api"
	"orgchart_nexoan/models"
)

func main() {
	// Create a new API client
	client := api.NewClient("http://localhost:8080/entities", "http://localhost:8081/v1/entities")

	// Example 1: Get root entities
	rootEntities, err := client.GetRootEntities("organization")
	if err != nil {
		log.Fatalf("Failed to get root entities: %v", err)
	}
	fmt.Printf("Root entities: %v\n", rootEntities)

	// Example 2: Search for entities
	searchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "organization",
			Minor: "company",
		},
		Name: "Example Corp",
	}
	searchResults, err := client.SearchEntities(searchCriteria)
	if err != nil {
		log.Fatalf("Failed to search entities: %v", err)
	}
	fmt.Printf("Search results: %+v\n", searchResults)

	// Example 3: Get entity metadata
	if len(searchResults) > 0 {
		metadata, err := client.GetEntityMetadata(searchResults[0].ID)
		if err != nil {
			log.Fatalf("Failed to get entity metadata: %v", err)
		}
		fmt.Printf("Entity metadata: %+v\n", metadata)
	}

	// Example 4: Get entity attribute
	if len(searchResults) > 0 {
		attribute, err := client.GetEntityAttribute(
			searchResults[0].ID,
			"location",
			time.Now().Format(time.RFC3339), // start time
			"",                              // end time (optional)
		)
		if err != nil {
			log.Fatalf("Failed to get entity attribute: %v", err)
		}
		fmt.Printf("Entity attribute: %+v\n", attribute)
	}

	// Example 5: Get related entities
	if len(searchResults) > 0 {
		query := &models.Relationship{
			RelatedEntityID: "some-other-entity",
			StartTime:       time.Now().Format(time.RFC3339),
			EndTime:         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			ID:              "relation-id",
			Name:            "relation-name",
		}
		relations, err := client.GetRelatedEntities(searchResults[0].ID, query)
		if err != nil {
			log.Fatalf("Failed to get related entities: %v", err)
		}
		fmt.Printf("Related entities: %+v\n", relations)
	}

	// Example 6: Get all related entities
	if len(searchResults) > 0 {
		allRelations, err := client.GetAllRelatedEntities(searchResults[0].ID)
		if err != nil {
			log.Fatalf("Failed to get all related entities: %v", err)
		}
		fmt.Printf("All related entities: %+v\n", allRelations)
	}
}
