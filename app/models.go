package app

type Example []interface{}

type Object struct {
	SchemaName  string          `json:"$schema,omitempty"`
	Id          string          `json:"$id"`
	Type        Type            `json:"type,omitempty"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Default     interface{}     `json:"default,omitempty"`
	Required    []string        `json:"required,omitempty"`
	Properties  map[string]temp `json:"properties,omitempty"`
	// Examples    *Example               `json:"examples,omitempty"`
	// Items       map[string]interface{} `json:"items,omitempty"`
}

type Property struct {
	*Object
}
