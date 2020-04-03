package models

import "encoding/json"

// Component defines the properties of a component
type Component struct {
	Name  string          `json:"component"`
	Props json.RawMessage `json:"props"`
}
