package types

import (
	"strings"
	"sync"
	"unicode"
)

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

func ToCamelCase(str string, initCase bool) string {
	separators := "-#@!$&=.+:;_~ (){}[]/"
	s := strings.Trim(str, " ")

	n := ""
	capNext := initCase
	for i, v := range s {
		if unicode.IsUpper(v) && i == 0 && !capNext {
			n += strings.ToLower(string(v))
		} else {
			if unicode.IsUpper(v) {
				n += string(v)
			}
			if unicode.IsDigit(v) {
				n += string(v)
			}
			if unicode.IsLower(v) {
				if capNext {
					n += strings.ToUpper(string(v))
				} else {
					n += string(v)
				}
			}

			if strings.ContainsRune(separators, v) {
				capNext = true
			} else {
				capNext = false
			}
		}
	}

	return n
}

// RemoveWhiteSpaceAndCaps This function is used to remove any white space from the provided string, and.. if white space is found, ensures
// the next character is capitalized. This helps camel case strings and should only be used for such requirements
func RemoveWhiteSpaceAndCaps(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	capNextCh := false

	for _, ch := range str {
		if unicode.IsSpace(ch) {
			capNextCh = true
		} else {
			if capNextCh {
				capNextCh = false
				b.WriteRune(unicode.ToUpper(ch))
			} else {
				b.WriteRune(ch)
			}
		}
	}

	return b.String()
}
