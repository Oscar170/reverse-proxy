package models

import "encoding/json"

// CompoentRendered defines the properties of a rendered component
type CompoentRendered struct {
	Html      string          `json:"html"`
	Css       string          `json:"css"`
	InitState json.RawMessage `json:"initState"`
}
