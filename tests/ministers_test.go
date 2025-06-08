package tests

import (
	"fmt"
	"orgchart_nexoan/api"
	"orgchart_nexoan/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *api.Client

func TestMain(m *testing.M) {
	// Set up test environment with correct URLs
	client = api.NewClient("http://localhost:8080/entities", "http://localhost:8081/v1/entities")

	// Create government node using CreateGovernmentNode
	government, err := client.CreateGovernmentNode()
	if err != nil {
		fmt.Printf("Failed to create government node: %v\n", err)
		os.Exit(1)
	}
	if government == nil {
		fmt.Println("Government node is nil")
		os.Exit(1)
	}
	fmt.Printf("Successfully created government node with ID: %s\n", government.ID)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateMinisters(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 0,
	}

	// Test cases for creating ministers
	testCases := []struct {
		transactionID string
		parent        string
		parentType    string
		child         string
		childType     string
		relType       string
		date          string
	}{
		{
			transactionID: "2153/12_tr_01",
			parent:        "Government of Sri Lanka",
			parentType:    "government",
			child:         "Minister of Defence",
			childType:     "minister",
			relType:       "AS_MINISTER",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_02",
			parent:        "Government of Sri Lanka",
			parentType:    "government",
			child:         "Minister of Finance, Economic and Policy Development",
			childType:     "minister",
			relType:       "AS_MINISTER",
			date:          "2019-12-10",
		},
	}

	// Create each minister
	for _, tc := range testCases {
		t.Logf("Creating minister: %s", tc.child)

		// Create transaction map for AddEntity
		transaction := map[string]interface{}{
			"parent":         tc.parent,
			"child":          tc.child,
			"date":           tc.date,
			"parent_type":    tc.parentType,
			"child_type":     tc.childType,
			"rel_type":       tc.relType,
			"transaction_id": tc.transactionID,
		}

		// Use AddEntity to create the minister
		_, err := client.AddEntity(transaction, entityCounters)
		assert.NoError(t, err)

		// Update the counter for the next iteration
		entityCounters[tc.childType]++

		// Verify the minister was created by searching for it
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.childType,
			},
			Name: tc.child,
		}

		results, err := client.SearchEntities(searchCriteria)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, tc.child, results[0].Name)

		// Verify the relationship was created by checking parent's relationships
		parentResults, err := client.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.parentType,
			},
			Name: tc.parent,
		})
		assert.NoError(t, err)
		assert.Len(t, parentResults, 1)

		// Get parent's metadata to verify relationship
		metadata, err := client.GetEntityMetadata(parentResults[0].ID)
		assert.NoError(t, err)
		assert.NotNil(t, metadata)
	}
}

func TestCreateDepartments(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"department": 0,
	}

	// Test cases for creating departments
	testCases := []struct {
		transactionID string
		parent        string
		parentType    string
		child         string
		childType     string
		relType       string
		date          string
	}{
		{
			transactionID: "2153/12_tr_03",
			parent:        "Minister of Defence",
			parentType:    "minister",
			child:         "Sri Lankan Army",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_04",
			parent:        "Minister of Finance, Economic and Policy Development",
			parentType:    "minister",
			child:         "Department of Taxes",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_05",
			parent:        "Minister of Finance, Economic and Policy Development",
			parentType:    "minister",
			child:         "Department of Policies",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
	}

	// Create each department
	for _, tc := range testCases {
		t.Logf("Creating department: %s under minister: %s", tc.child, tc.parent)

		// Create transaction map for AddEntity
		transaction := map[string]interface{}{
			"parent":         tc.parent,
			"child":          tc.child,
			"date":           tc.date,
			"parent_type":    tc.parentType,
			"child_type":     tc.childType,
			"rel_type":       tc.relType,
			"transaction_id": tc.transactionID,
		}

		// Use AddEntity to create the department
		_, err := client.AddEntity(transaction, entityCounters)
		assert.NoError(t, err)

		// Update the counter for the next iteration
		entityCounters[tc.childType]++

		// Verify the department was created by searching for it
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.childType,
			},
			Name: tc.child,
		}

		results, err := client.SearchEntities(searchCriteria)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, tc.child, results[0].Name)

		// Verify the relationship was created by checking minister's relationships
		ministerResults, err := client.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.parentType,
			},
			Name: tc.parent,
		})
		assert.NoError(t, err)
		assert.Len(t, ministerResults, 1)

		// Get minister's relationships to verify department relationship
		// relations, err := client.GetRelatedEntities(ministerResults[0].ID, &models.Relationship{
		// 	Name:    tc.relType,
		// 	EndTime: "",
		// })
		// assert.NoError(t, err)
		// assert.NotEmpty(t, relations, "Minister should have at least one department relationship")
	}
}

func TestTerminateDepartment(t *testing.T) {
	// Create transaction map for terminating the department
	transaction := map[string]interface{}{
		"parent":      "Minister of Defence",
		"child":       "Sri Lankan Army",
		"date":        "2024-01-01",
		"parent_type": "minister",
		"child_type":  "department",
		"rel_type":    "AS_DEPARTMENT",
	}

	// Terminate the department relationship
	err := client.TerminateEntity(transaction)
	assert.NoError(t, err)

	// Find the minister to verify the relationship
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Defence",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	// Find the department
	departmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Sri Lankan Army",
	})
	assert.NoError(t, err)
	assert.Len(t, departmentResults, 1)
	departmentID := departmentResults[0].ID

	// Verify the relationship is terminated
	allRelations, err := client.GetAllRelatedEntities(ministerID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated relationship")
}

func TestTerminateMinister(t *testing.T) {
	// Create transaction map for terminating the minister
	transaction := map[string]interface{}{
		"parent":      "Government of Sri Lanka",
		"child":       "Minister of Defence",
		"date":        "2024-01-01",
		"parent_type": "government",
		"child_type":  "minister",
		"rel_type":    "AS_MINISTER",
	}

	// Terminate the minister relationship
	err := client.TerminateEntity(transaction)
	assert.NoError(t, err)

	// Find the government to verify the relationship
	governmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: "Government of Sri Lanka",
	})
	assert.NoError(t, err)
	assert.Len(t, governmentResults, 1)
	governmentID := governmentResults[0].ID

	// Find the minister
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Defence",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	// Verify the relationship is terminated
	allRelations, err := client.GetAllRelatedEntities(governmentID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == ministerID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated relationship")
}

func TestMoveDepartment(t *testing.T) {
	// Create transaction map for moving the department
	transaction := map[string]interface{}{
		"old_parent": "Minister of Finance, Economic and Policy Development",
		"new_parent": "Minister of Defence",
		"child":      "Department of Policies",
		"type":       "AS_DEPARTMENT",
		"date":       "2024-01-01",
	}

	// Move the department
	err := client.MoveDepartment(transaction)
	assert.NoError(t, err)

	// Find the new minister to verify the new relationship
	newMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Defence",
	})
	assert.NoError(t, err)
	assert.Len(t, newMinisterResults, 1)
	newMinisterID := newMinisterResults[0].ID

	// Find the department
	departmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Department of Policies",
	})
	assert.NoError(t, err)
	assert.Len(t, departmentResults, 1)
	departmentID := departmentResults[0].ID

	// Verify the new relationship exists
	allRelations, err := client.GetAllRelatedEntities(newMinisterID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the new relationship")

	// Find the old minister to verify the old relationship is terminated
	oldMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance, Economic and Policy Development",
	})
	assert.NoError(t, err)
	assert.Len(t, oldMinisterResults, 1)
	oldMinisterID := oldMinisterResults[0].ID

	// Verify the old relationship is terminated
	oldRelations, err := client.GetAllRelatedEntities(oldMinisterID)
	assert.NoError(t, err)
	found = false
	for _, rel := range oldRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated old relationship")
}

// TODO: Test that it fails and returns proper error messages when trying to terminate an entity with children
