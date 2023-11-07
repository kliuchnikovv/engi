package parameter

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type body struct {
	pointer   interface{}
	unmarshal func(*request.Request) (types.Unmarshaler, error)
	options   []request.Option
}

func (p *body) Handle(r *request.Request, w http.ResponseWriter) error {
	unmarshal, err := p.unmarshal(r)
	if err != nil {
		return err
	}

	return request.ExtractBody(r, unmarshal, p.pointer, p.options)
}

// Body - takes pointer to structure and saves casted request body into context and pointer.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func Body(pointer interface{}, opts ...request.Option) request.Middleware {
	return &body{
		pointer: pointer,
		options: opts,
		unmarshal: request.GetUnmarshaler,
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func CustomBody(unmarshal types.Unmarshaler, pointer interface{}, opts ...request.Option) request.Middleware {
	return &body{
		pointer: pointer,
		options: opts,
		unmarshal: func(*request.Request) (types.Unmarshaler, error) {
			return unmarshal, nil
		},
	}
}
