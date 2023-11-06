package options

import "github.com/KlyuchnikovV/engi/internal/request"

// Description - adds string description to parameter.
// Can be used in errors description or in documentation.
func Description(s string) request.Option {
	return func(p *request.Parameter) error {
		p.Description = s

		return nil
	}
}
