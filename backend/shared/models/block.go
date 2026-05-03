package models

// Block represents a recursive component block for UI construction
type Block struct {
	ID           string                 `bson:"id" json:"id"`
	ComponentKey string                 `bson:"component_key" json:"componentKey"`
	Props        map[string]interface{} `bson:"props,omitempty" json:"props,omitempty"`
	Children     []Block                `bson:"children,omitempty" json:"children,omitempty"`
}
