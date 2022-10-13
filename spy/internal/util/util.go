package util

import (
	"fmt"
	"time"
)

func FileNameFormat(name, ext string) string {
	return fmt.Sprintf("%s-%s%s", name, time.Now().Format("2006-01-02-15-04-05"), ext)
}

func Now() string {
	return time.Now().Format(time.RFC3339)
}
