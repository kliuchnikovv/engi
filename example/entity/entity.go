package entity

type NotesRequest struct {
	Note   string `json:"note" description:"Note content in Markdown" example:"# Heading level 1"`
	Author string `json:"author" description:"Author name" example:"John Cane"`
}

type RequestBody struct {
	String       string      `json:"field" description:"Just a string"`
	Integer      int         `json:"integer,omitempty"`
	SimpleArray  []string    `json:"simpleArray"`
	ArrayOfArray [][]float32 `json:"arrayOfArray"`
	WithoutTag   float64
}
