package dbg

import (
	"fmt"

	w32 "github.com/mjwhitta/win/api"
)

// Printf is a wrapper for OutputDebugStringW that supports format
// strings.
func Printf(format string, params ...any) {
	w32.OutputDebugStringW(fmt.Sprintf(format, params...))
}
