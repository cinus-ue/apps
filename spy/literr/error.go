package literr

import (
	"errors"
	"log"
	"os"
)

var ArgsError = errors.New("required arguments not provided")

func CheckError(errs ...error) bool {
	if !Discard {
		LogError(errs...)
	}
	return HasError(errs...)
}

func CheckFatal(errs ...error) {
	if CheckError(errs...) {
		os.Exit(1)
	}
}
func LogError(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		log.Println(" Error: ", err.Error())
		if Debug {
			LogStack()
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
