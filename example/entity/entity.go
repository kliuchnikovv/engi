package entity

type NotesRequest struct {
	Note   string
	Author string
}

type RequestBody struct {
	String       string      `json:"field"`
	Integer      int         `json:"integer"`
	SimpleArray  []string    `json:"simpleArray"`
	ArrayOfArray [][]float32 `json:"arrayOfArray"`
	WithoutTag   float64
}
