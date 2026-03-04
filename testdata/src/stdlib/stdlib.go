package stdlib

import (
	"context"
	"io"
)

// StdlibOK demonstrates that stdlib interfaces are always allowed.
type StdlibOK struct {
	reader io.Reader       // OK: stdlib interface
	writer io.Writer       // OK: stdlib interface
	closer io.Closer       // OK: stdlib interface
	ctx    context.Context // OK: stdlib interface
}
