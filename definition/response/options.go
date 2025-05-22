package response

import (
	"github.com/kliuchnikovv/engi/internal/routes"
	"github.com/kliuchnikovv/engi/internal/types"
)

func ResponseAs(responser func() Responser) routes.Middleware {
	return &responserObject{
		responser: responser(),
	}
}

func MarshalAs(marshaler func() Marshaler) routes.Middleware {
	return &marshalerObject{
		marshaler: types.Marshaler(marshaler()),
	}
}

func AsIs() Responser {
	return new(types.ResponseAsIs)
}

func AsObject() Responser {
	return new(types.ResponseAsObject)
}

func AsJSON() Marshaler {
	return Marshaler(types.NewJSONMarshaler())
}

func AsXML() Marshaler {
	return Marshaler(types.NewXMLMarshaler())
}
