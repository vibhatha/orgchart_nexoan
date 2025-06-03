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
