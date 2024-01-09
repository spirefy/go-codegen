package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var pathParamRE *regexp.Regexp

type Formatter interface {
	Format([]byte) ([]byte, error)
}

func init() {
	pathParamRE = regexp.MustCompile("{[.;?]?([^{}*]+)\\*?}")
}

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

// Debug
// PrintObj Helper "debug" function to print out the JSON equivalent structure of a Go object
// numLines indicates how many lines to print before and/or after the object to help separate it from
// other output. before and after are true to print the output to numLines before and/or after the
// object is printed.
func Debug(obj interface{}, numLines int, before, after bool) {
	jsonData, _ := json.MarshalIndent(obj, "", "    ")
	if before {
		for i := 1; i < numLines; i++ {
			CLog.Log(CLogCore, "", "**********************************************************************")
		}
	}

	CLog.Log(CLogCore, "", "")

	CLog.Log(CLogCore, "", "%s", jsonData)
	CLog.Log(CLogCore, "", "")
	if after {
		for i := 1; i < numLines; i++ {
			CLog.Log(CLogCore, "", "**********************************************************************")
		}
	}
}

func GetRootPath() *string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil
	}

	return &dir
}

// WriteToFile
//
// This helper function will accept the path and filename to write out as a file, but
// as this is GO code being passed in the contents slice.. it will use the Go std format() (gofmt)
// to format the generated code before writing it out.
func WriteToFile(path string, filename string, contents []byte, formatter Formatter, logCategory CodegenLogCategory) error {
	if len(path) > 0 && len(filename) > 0 && contents != nil {
		// make sure dirs exist.. or create them:
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			CLog.Error(logCategory, "", "%", err)
			return err
		}

		f, err := os.Create(path + string(os.PathSeparator) + filename)

		if err != nil {
			CLog.Error(logCategory, "", "%S", err)
			return err
		}

		defer func(f *os.File) {
			err := f.Close()

			if err != nil {
				CLog.Error(logCategory, "", "Error trying to defer closing of file. %S", err)
			}
		}(f)

		// Format code using language formatter if provided
		formatted := contents

		if nil != formatter {
			formatted, err = formatter.Format(contents)

			if err != nil {
				CLog.Error(logCategory, "", "Error trying to format code: ", err)
			}
		}

		_, err = f.Write(formatted)

		if err != nil {
			CLog.Error(logCategory, "", "%s", err)
			return err
		}
	}

	return nil
}

func CamelCase(str string) string {
	return ToCamelCase(str, false)
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

// UriPathParamToColonParam This function converts a swagger style path URI with parameters to a
// Echo compatible path URI. We need to replace all of Swagger parameters with
// ":param". Valid input parameters are:
//
//	{param}
//	{param*}
//	{.param}
//	{.param*}
//	{;param}
//	{;param*}
//	{?param}
//	{?param*}
func UriPathParamToColonParam(uri string) string {
	return pathParamRE.ReplaceAllString(uri, ":$1")
}

// UriPathParamToBracesParam This function converts a swagger style path URI with parameters to a
// Chi compatible path URI. We need to replace all of Swagger parameters with
// "{param}". Valid input parameters are:
//
//	{param}
//	{param*}
//	{.param}
//	{.param*}
//	{;param}
//	{;param*}
//	{?param}
//	{?param*}
func UriPathParamToBracesParam(uri string) string {
	return pathParamRE.ReplaceAllString(uri, "{$1}")
}

// UriPathParamToAngleBracketsParam This function converts a swagger style path URI with parameters to a
// Flask compatible path URI. We need to replace all Swagger parameters with
// <param>
func UriPathParamToAngleBracketsParam(uri string) string {
	return pathParamRE.ReplaceAllStringFunc(uri, func(param string) string {
		param = strings.ReplaceAll(param, "-", "_")
		param = strings.ReplaceAll(param, "{", "<")
		param = strings.ReplaceAll(param, "}", ">")
		return param
	})
}

// UriPathParamToLowerCamelParam This function converts swagger style path to /:userId/house/:houseId removes the hyphens
// from the path params and makes it camelcase.
func UriPathParamToLowerCamelParam(uri string) string {
	return pathParamRE.ReplaceAllStringFunc(uri, func(param string) string {
		param = strings.ReplaceAll(param, "{", "")
		param = strings.ReplaceAll(param, "}", "")
		param = ":" + ToCamelCase(param, false)
		return param
	})
}

// ToLowerCamel This converts "AnyKind of_string" to "anyKindOfString"
func ToLowerCamel(s string) string {
	return ToCamelCase(s, false)
}

// StripNonAlphaNumeric This function will remove any non-alpha characters from the provided string
func StripNonAlphaNumeric(s string) string {
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}

	return result.String()
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

func Camelize(str string) string {
	var g []string
	p := strings.Fields(str)

	for _, value := range p {
		g = append(g, strings.Title(value))
	}

	return strings.Join(g, "")
}

// ToUnderscoreCase This function will convert query-arg style strings to Underscore Case. We will
// use `., -, +, :, ;, _, ~, ' ', (, ), {, }, [, ]` as valid delimiters for words.
// So, "word.word-word+word:word;word_word~word word(word)word{word}[word]"
// would be converted to word_word_word_word_word_word_word_word_word_word_word_word_word
func ToUnderscoreCase(str string) string {
	separators := "-#@!$&=.+:;_~ (){}[]"
	s := strings.Trim(str, " ")
	n := ""
	// Add an underscore if the string starts with a number
	if s[0] >= '0' && s[0] <= '9' {
		n += "_"
	}

	for _, v := range s {
		if strings.ContainsRune(separators, v) {
			n += "_"
		} else {
			n += string(v)
		}
	}
	return n
}

// StringArrayContains Contains tells whether a contains x.
func StringArrayContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}

func ExtensionNameValue(extPropValue interface{}) (string, error) {
	raw, ok := extPropValue.(json.RawMessage)

	if !ok {
		return "", fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return "", errors.Errorf("a random error because of %s\n", err)
	}

	return name, nil
}

// TypeToString This function will attempt to convert the value provided to a String.. and if necessary use the format provided
// as a guide
func TypeToString(value interface{}, format string) string {
	tok, ok := value.(string)

	if !ok && format != "" {
		switch format {
		case "float64":
			tok = fmt.Sprintf("%f", value.(float64))
		case "float32":
			tok = fmt.Sprintf("%f", value.(float32))
		case "int64":
			tok = strconv.FormatInt(value.(int64), 10)
		case "int32":
			tok = string(value.(int32))
		case "int16":
			tok = strconv.Itoa(int(value.(int16)))
		case "int8":
			tok = strconv.Itoa(int(value.(int8)))
		case "uint":
			tok = strconv.Itoa(int(value.(uint)))
		case "int":
			tok = string(rune(value.(int)))
		case "string":
			tok = value.(string)
		case "uint16":
			tok = strconv.Itoa(int(value.(uint16)))
		case "uint32":
			tok = string(value.(int32))
		case "uint64":
			tok = strconv.FormatInt(value.(int64), 10)
		case "bool":
			tok = strconv.FormatBool(value.(bool))
		}
	}

	return tok
}
