package models

import (
	"encoding/base64"
	"encoding/json"
)

// SearchCriteria represents the search parameters for entity search
type SearchCriteria struct {
	Kind       *Kind  `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Created    string `json:"created,omitempty"`
	Terminated string `json:"terminated,omitempty"`
}

// SearchResult represents a single entity in search results
type SearchResult struct {
	ID         string `json:"id"`
	Kind       Kind   `json:"kind"`
	Name       string `json:"name"`
	Created    string `json:"created"`
	Terminated string `json:"terminated,omitempty"`
}

// SearchResponse represents the response from the search endpoint
type SearchResponse struct {
	Body []SearchResult `json:"body"`
}

// RootEntitiesResponse represents the response from the root entities endpoint
type RootEntitiesResponse struct {
	Body []string `json:"body"`
}

// AttributeValue represents a single time-based attribute value
type AttributeValue struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
	Value string `json:"value"`
}

// RelationQuery represents the query parameters for getting related entities
type RelationQuery struct {
	RelatedEntityID string `json:"relatedEntityId"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	ID              string `json:"id"`
	Name            string `json:"name"`
}

// Relation represents a relationship between entities
type Relation struct {
	RelatedEntityID string `json:"relatedEntityId"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	ID              string `json:"id"`
	Name            string `json:"name"`
}

// UnmarshalJSON implements custom JSON unmarshaling for SearchResult
func (s *SearchResult) UnmarshalJSON(data []byte) error {
	type Alias SearchResult
	aux := &struct {
		Name json.RawMessage `json:"name"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Try to unmarshal as a simple string first
	var simpleName string
	if err := json.Unmarshal(aux.Name, &simpleName); err == nil {
		s.Name = simpleName
		return nil
	}

	// If that fails, try to unmarshal as a protobuf string value
	var protobufName struct {
		TypeURL string `json:"typeUrl"`
		Value   string `json:"value"`
	}
	if err := json.Unmarshal(aux.Name, &protobufName); err != nil {
		return err
	}

	// Decode the base64 value
	decoded, err := base64.StdEncoding.DecodeString(protobufName.Value)
	if err != nil {
		return err
	}
	s.Name = string(decoded)
	return nil
}
