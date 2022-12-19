package util

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cinus-ue/spy/literr"
)

const (
	CharEscape      = '\\'
	CharSingleQuote = '\''
	CharDoubleQuote = '"'
	CharBackQuote   = '`'
)

func IsQuote(r rune) bool {
	return r == CharSingleQuote || r == CharDoubleQuote || r == CharBackQuote
}

// ParseArgs parses line, ignore brackets
func ParseArgs(line string) (lineArgs []string) {
	var (
		rl        = []rune(line + " ")
		buf       = strings.Builder{}
		quoteChar rune
		nextChar  rune
		escaped   bool
		in        bool
	)

	var isSpace bool

	for k, r := range rl {
		isSpace = unicode.IsSpace(r)
		if !isSpace && !in {
			in = true
		}
		switch {
		case escaped:
			escaped = false
			//pass
		case r == CharEscape: // Escape mode
			if k+1+1 < len(rl) {
				nextChar = rl[k+1]
				// Only these characters are supported for escaping,
				// otherwise the backslash is output as-is
				if unicode.IsSpace(nextChar) || IsQuote(nextChar) || nextChar == CharEscape {
					escaped = true
					continue
				}
			}
			// pass
		case IsQuote(r):
			if quoteChar == 0 {
				quoteChar = r
				continue
			}

			if quoteChar == r {
				quoteChar = 0
				continue
			}
		case isSpace:
			if !in { // ignore space
				continue
			}
			if quoteChar == 0 { // Not in quotes
				lineArgs = append(lineArgs, buf.String())
				buf.Reset()
				in = false
				continue
			}
		}
		buf.WriteRune(r)
	}
	return
}

func FileNameFormat(name, ext string) string {
	return fmt.Sprintf("%s-%s%s", name, time.Now().Format("2006-01-02-15-04-05"), ext)
}

func Now() string {
	return time.Now().Format(time.RFC3339)
}

func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

func ParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	literr.CheckFatal(err)
	return t
}

func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	literr.CheckError(err)
	return i
}

func Unique(s []string) []string {
	m := make(map[string]bool)
	var ret []string
	for _, i := range s {
		if _, ok := m[i]; !ok {
			m[i] = true
			ret = append(ret, i)
		}
	}
	return ret
}

func HttpClient() *http.Client {
	tr := &http.Transport{}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	hc := &http.Client{Transport: tr}
	return hc
}
