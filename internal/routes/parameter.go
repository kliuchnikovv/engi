package routes

type NamedParameter interface {
	Name() string
	Regexp() string
}

// type IParameter interface {
// 	Bind(*request.Request, *response.Response) error
// }

// type ParameterType interface {
// 	bool | int64
// }

// type ParameterHandler func(request *request.Request, response *response.Response) error

// type Parameter[T ParameterType] struct {
// 	key     string
// 	place   placing.Placing
// 	options []request.Option

// 	handler ParameterHandler
// 	value   T
// }

// func NewParameter(
// 	key string,
// 	place placing.Placing,
// 	handler ParameterHandler,
// ) *Parameter {
// 	return &Parameter{
// 		key:     key,
// 		place:   place,
// 		handler: handler,
// 	}
// }

// func (parameter Parameter[T]) Regexp() {

// }

// func (parameter Parameter[T]) BindValue(
// 	request *request.Request,
// 	response *response.Response,
// 	convert func(string) (T, error),
// ) error {
// 	var rawParam = request.GetParameter(parameter.key, parameter.place)
// 	if len(rawParam) == 0 {
// 		return fmt.Errorf("parameter not found: '%s'", parameter.key)
// 	}

// 	// result, err := convert(rawParam)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// if result != nil {
// 	// 	var parameter = request.parameters[paramPlacing][key]

// 	// 	request.parameters[paramPlacing][key] = Parameter{
// 	// 		Name:         key,
// 	// 		Parsed:       result,
// 	// 		raw:          parameter.raw,
// 	// 		Description:  parameter.Description,
// 	// 		wasRequested: true,
// 	// 	}
// 	// }

// 	// var parameter = request.parameters[paramPlacing][key]
// 	// for _, config := range configs {
// 	// 	if err := config(&parameter); err != nil {
// 	// 		return err
// 	// 	}
// 	// }

// 	// parameter.Name = key
// 	// request.parameters[paramPlacing][key] = parameter

// 	return nil

// 	// return strconv.ParseBool(request)
// }

// func (param Parameter[T]) Value() T {
// 	return param.value
// }

// type Boolean Parameter[bool]

// func (parameter *Boolean) Bind(
// 	request *request.Request,
// 	response *response.Response,
// ) error {
// 	var raw = request.GetParameter(parameter.key, parameter.place)
// 	if len(raw) == 0 {
// 		return fmt.Errorf("parameter not found: '%s'", parameter.key)
// 	}

// 	value, err := strconv.ParseBool(raw)
// 	if err != nil {
// 		return err
// 	}

// 	return request.UpdateParameter(
// 		response,
// 		parameter.key,
// 		parameter.place,
// 		value,
// 		parameter.options...,
// 	)
// }
