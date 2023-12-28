package errors

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidationErrors is a list of field errors.
// Avoid using this type directly, use ValidationErrorBuilder instead.
// Because, you need to add checks like `if len(ve) == 0 return nil`
// Which is error prone, so use ValidationErrorBuilder instead.
type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	return "validation failed"
}

func ValidationErrs() *ValidationErrorBuilder {
	return &ValidationErrorBuilder{}
}

type ValidationErrorBuilder struct {
	ve ValidationErrors
}

func (b *ValidationErrorBuilder) Add(field, err string) {
	b.ve = append(b.ve, FieldError{Field: field, Error: err})
}

func (b *ValidationErrorBuilder) Err() error {
	if len(b.ve) == 0 {
		return nil
	}
	return b.ve
}
