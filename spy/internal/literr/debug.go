package literr

import "runtime/debug"

var (
	Discard bool
	Debug   bool
)

func PrintStack() {
	debug.PrintStack()
}
