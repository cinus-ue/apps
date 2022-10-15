package literr

import (
	"errors"
	"fmt"
	"os"
)

var ArgsError = errors.New("required arguments not provided")

func CheckError(errs ...error) bool {
	if !Discard {
		PrintError(errs...)
	}
	return HasError(errs...)
}

func CheckFatal(errs ...error) {
	if CheckError(errs...) {
		os.Exit(1)
	}
}
func PrintError(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		fmt.Fprintln(os.Stderr, " Error: ", err.Error())
		if Debug {
			PrintStack()
		}
	}
}

func HasError(errs ...error) bool {
	if len(errs) == 0 {
		return false
	}
	hasError := false
	for _, err := range errs {
		if err == nil {
			continue
		}
		hasError = true
	}
	return hasError
}
