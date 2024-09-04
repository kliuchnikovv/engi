package response

import (
	"github.com/KlyuchnikovV/engi/internal/routes"
	"github.com/KlyuchnikovV/engi/internal/types"
)

func ResponseAs(responser func() Responser) routes.Option {
	return &responserObject{
		responser: responser(),
	}
}

func MarshalAs(marshaler func() Marshaler) routes.Option {
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
