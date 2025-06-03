package models

// Entity represents the main entity structure
type Entity struct {
	ID            string              `json:"id,omitempty"`
	Kind          Kind                `json:"kind,omitempty"`
	Created       string              `json:"created,omitempty"`
	Terminated    string              `json:"terminated,omitempty"`
	Name          TimeBasedValue      `json:"name"`
	Metadata      []MetadataEntry     `json:"metadata,omitempty"`
	Attributes    []AttributeEntry    `json:"attributes,omitempty"`
	Relationships []RelationshipEntry `json:"relationships,omitempty"`
}

// Kind represents the entity kind structure
type Kind struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
}

// TimeBasedValue represents a value that changes over time
type TimeBasedValue struct {
	StartTime string      `json:"startTime"`
	EndTime   string      `json:"endTime,omitempty"`
	Value     interface{} `json:"value"`
}

// AttributeEntry represents a key-value pair in the attributes array
type AttributeEntry struct {
	Key   string                   `json:"key"`
	Value AttributeValueCollection `json:"value"`
}

// AttributeValueCollection represents the value structure for attributes
type AttributeValueCollection struct {
	Values []TimeBasedValue `json:"values"`
}

// MetadataEntry represents a key-value pair in the metadata array
type MetadataEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// RelationshipEntry represents a key-value pair in the relationships array
type RelationshipEntry struct {
	Key   string       `json:"key"`
	Value Relationship `json:"value"`
}

// Relationship represents a relationship between entities
type Relationship struct {
	RelatedEntityID string `json:"relatedEntityId"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime,omitempty"`
	ID              string `json:"id"`
	Name            string `json:"name"`
}
