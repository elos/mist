package middleware

import "github.com/elos/ehttp/middleware"

// Cors wraps the ehttp/middleware.Cors type
type Cors struct {
	*middleware.Cors
}

// NewCors constructs a new Cors object with Allowed Headers = headers
func NewCors(headers ...string) *Cors {
	return &Cors{middleware.NewCors(headers...)}
}
