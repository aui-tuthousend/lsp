package http

import "github.com/danielgtaylor/huma/v2"

func WithDoc(operationID, summary, description string, tags ...string) func(*huma.Operation) {
	return func(o *huma.Operation) {
		o.OperationID = operationID
		o.Summary = summary
		o.Description = description
		o.Tags = tags
	}
}

// WithStatus menetapkan default HTTP status code untuk sebuah operasi
// (mis. 201 Created).
func WithStatus(code int) func(*huma.Operation) {
	return func(o *huma.Operation) {
		o.DefaultStatus = code
	}
}
