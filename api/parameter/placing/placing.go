package placing

type Placing string

const (
	InPath   Placing = "path"
	InQuery  Placing = "query"
	InCookie Placing = "cookie"
	InHeader Placing = "header"
)
