package api

import (
	"fmt"
	"strings"
	"time"

	"orgchart_nexoan/models"
)

// CreateGovernmentNode creates the initial government node
func (c *Client) CreateGovernmentNode() (*models.Entity, error) {
	// Create the government entity
	governmentEntity := &models.Entity{
		ID:      "gov_01",
		Created: "2024-01-01T00:00:00Z",
		Kind: models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: models.TimeBasedValue{
			StartTime: "2024-01-01T00:00:00Z",
			Value:     "Government of Sri Lanka",
		},
	}

	// Create the entity
	createdEntity, err := c.CreateEntity(governmentEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create government entity: %w", err)
	}

	return createdEntity, nil
}

// AddEntity creates a new entity and establishes its relationship with a parent entity.
// Assumes the parent entity already exists.
func (c *Client) AddEntity(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)
	transactionID := transaction["transaction_id"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Generate new entity ID
	if _, exists := entityCounters[childType]; !exists {
		return 0, fmt.Errorf("unknown child type: %s", childType)
	}

	prefix := fmt.Sprintf("%s_%s", transactionID[:7], strings.ToLower(childType[:3]))
	entityCounter := entityCounters[childType] + 1
	newEntityID := fmt.Sprintf("%s_%d", prefix, entityCounter)

	// Get the parent entity ID
	searchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: parentType,
		},
		Name: parent,
	}

	fmt.Printf("[AddEntity] Searching for parent with criteria: %+v\n", searchCriteria)
	searchResults, err := c.SearchEntities(searchCriteria)
	fmt.Printf("[AddEntity] Search Results: %+v\n", searchResults)
	if err != nil {
		return 0, fmt.Errorf("failed to search for parent entity: %w", err)
	}

	if len(searchResults) == 0 {
		return 0, fmt.Errorf("parent entity not found: %s", parent)
	}

	parentID := searchResults[0].ID
	fmt.Println("Parent ID: ", parentID)

	// Create the new child entity
	childEntity := &models.Entity{
		ID: newEntityID,
		Kind: models.Kind{
			Major: "Organisation",
			Minor: childType,
		},
		Created:    dateISO,
		Terminated: "",
		Name: models.TimeBasedValue{
			StartTime: dateISO,
			Value:     child,
		},
		Metadata:      []models.MetadataEntry{},
		Attributes:    []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{},
	}

	// Create the child entity
	createdChild, err := c.CreateEntity(childEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to create child entity: %w", err)
	}

	// Update the parent entity to add the relationship to the child
	parentEntity := &models.Entity{
		ID:         parentID,
		Kind:       models.Kind{},
		Created:    "",
		Terminated: "",
		Name:       models.TimeBasedValue{},
		Metadata:   []models.MetadataEntry{},
		Attributes: []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{
			{
				Key: fmt.Sprintf("%s_%s", parentID, createdChild.ID),
				Value: models.Relationship{
					RelatedEntityID: createdChild.ID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              fmt.Sprintf("%s_%s", parentID, createdChild.ID),
					Name:            relType,
				},
			},
		},
	}

	_, err = c.UpdateEntity(parentID, parentEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update parent entity: %w", err)
	}

	return entityCounter, nil
}
