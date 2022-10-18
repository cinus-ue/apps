package literr

import (
	"log"
	"runtime/debug"
)

var (
	Discard bool
	Debug   bool
)

func LogStack() {
	log.Println(debug.Stack())
}
