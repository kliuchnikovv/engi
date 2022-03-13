package webapi


// TODO: alternative registering

// type (
// 	Handler func(*Context) error

// 	Group interface {
// 		GET(string, ...Handler)
// 		POST(string, ...Handler)
// 		DELETE(string, ...Handler)
// 		PUT(string, ...Handler)
// 		HEAD(string, ...Handler)
// 		OPTIONS(string, ...Handler)
// 		PATCH(string, ...Handler)

// 		GroupAPI(string, ...Handler) Group
// 	}

// 	Engine interface {
// 		Group
// 		// Start(string) error
// 	}

// 	WebAPI struct {
// 		Engine
// 	}

// 	Context struct {
// 		Query    types.Query
// 		Response types.Response
// 	}
// )

// func NewContext(ctx types.APIContext) *Context {
// 	return &Context{
// 		Query:    query.New(ctx),
// 		Response: response.New(ctx),
// 	}
// }

// func New(e Engine) *WebAPI {
// 	return &WebAPI{
// 		Engine: e,
// 	}
// }

// type WebAPI struct {
// }

// func New() {
// 	// http.
// }
