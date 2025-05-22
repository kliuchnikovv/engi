package response

import "github.com/kliuchnikovv/engi/internal/types"

func SetMarshaler(resp *Response, marshaler types.Marshaler) {
	resp.marshaler = marshaler
}

func SetResponser(resp *Response, responser types.Responser) {
	resp.object = responser
}
