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

// TerminateEntity terminates a specific relationship between parent and child at a given date
func (c *Client) TerminateEntity(transaction map[string]interface{}) error {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the parent entity ID
	searchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: parentType,
		},
		Name: parent,
	}
	parentResults, err := c.SearchEntities(searchCriteria)
	if err != nil {
		return fmt.Errorf("failed to search for parent entity: %w", err)
	}
	if len(parentResults) == 0 {
		return fmt.Errorf("parent entity not found: %s", parent)
	}
	parentID := parentResults[0].ID

	// Get the child entity ID
	searchCriteria.Kind.Minor = childType
	searchCriteria.Name = child
	childResults, err := c.SearchEntities(searchCriteria)
	if err != nil {
		return fmt.Errorf("failed to search for child entity: %w", err)
	}
	if len(childResults) == 0 {
		return fmt.Errorf("child entity not found: %s", child)
	}
	childID := childResults[0].ID

	// Get the specific relationship that is still active (no end date) -> this should give us the relationship(s) active for dateISO
	relations, err := c.GetRelatedEntities(parentID, &models.Relationship{
		RelatedEntityID: childID,
		Name:            relType,
		StartTime:       dateISO,
	})
	if err != nil {
		return fmt.Errorf("failed to get relationship: %w", err)
	}

	fmt.Printf("[TerminateEntity] Relationships: for %s and %s with type %s: %+v\n", parentID, childID, relType, relations)

	// FIXME: Is it possible to have more than one active relationshoip? For orgchart case only it won't happen
	// Find the active relationship (no end time)
	var activeRel *models.Relationship
	for _, rel := range relations {
		if rel.EndTime == "" {
			activeRel = &rel
			break
		}
	}

	if activeRel == nil {
		return fmt.Errorf("no active relationship found between %s and %s with type %s", parentID, childID, relType)
	}

	// Update the relationship to set the end date
	_, err = c.UpdateEntity(parentID, &models.Entity{
		ID: parentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: activeRel.ID,
				Value: models.Relationship{
					EndTime: dateISO,
					ID:      activeRel.ID,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate relationship: %w", err)
	}

	return nil
}

// MoveDepartment moves a department from one minister to another
func (c *Client) MoveDepartment(transaction map[string]interface{}) error {
	// Extract details from the transaction
	newParent := transaction["new_parent"].(string)
	oldParent := transaction["old_parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	relType := transaction["type"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the new minister (parent) entity ID
	newParentResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: newParent,
	})
	if err != nil {
		return fmt.Errorf("failed to search for new parent entity: %w", err)
	}
	if len(newParentResults) == 0 {
		return fmt.Errorf("new parent entity not found: %s", newParent)
	}
	newParentID := newParentResults[0].ID

	// Get the department (child) entity ID
	childResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: child,
	})
	if err != nil {
		return fmt.Errorf("failed to search for child entity: %w", err)
	}
	if len(childResults) == 0 {
		return fmt.Errorf("child entity not found: %s", child)
	}
	childID := childResults[0].ID

	// Create new relationship between new minister and department
	newRelationship := &models.Entity{
		ID: newParentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: fmt.Sprintf("%s_%s", newParentID, childID),
				Value: models.Relationship{
					RelatedEntityID: childID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              fmt.Sprintf("%s_%s", newParentID, childID),
					Name:            relType,
				},
			},
		},
	}

	_, err = c.UpdateEntity(newParentID, newRelationship)
	if err != nil {
		return fmt.Errorf("failed to create new relationship: %w", err)
	}

	// Terminate the old relationship
	terminateTransaction := map[string]interface{}{
		"parent":      oldParent,
		"child":       child,
		"date":        dateStr,
		"parent_type": "minister",
		"child_type":  "department",
		"rel_type":    relType,
	}

	err = c.TerminateEntity(terminateTransaction)
	if err != nil {
		return fmt.Errorf("failed to terminate old relationship: %w", err)
	}

	return nil
}
