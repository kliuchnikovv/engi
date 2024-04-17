package entity

// TODO: add field checks
type NotesRequest struct {
	Note   string `description:"Note content in Markdown" example:"# Heading level 1" json:"note"`
	Author string `description:"Author name"              example:"John Cane"         json:"author"`
}

type RequestBody struct {
	String       string      `description:"Just a string" json:"field"`
	Integer      int         `json:"integer,omitempty"`
	SimpleArray  []string    `json:"simpleArray"`
	ArrayOfArray [][]float32 `json:"arrayOfArray"`
	WithoutTag   float64
}
