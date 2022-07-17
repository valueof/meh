package formatters

// Formatter is an interface implemented by types that can
// transform Medium export data into a different format.
type Formatter interface {
	WriteFile(fp string, v any) error
}
