package core

import (
	"fmt"
	"log"
	"strings"
)

// CodegenLogger
// structure to containt he map of allowed logging categories, and if it should be prefaced with Codegen Logger:
// to help distinguish console output from other loggin that may happen.
type CodegenLogger struct {
	categories        []CodegenLogCategory
	allowedCategories map[CodegenLogCategory]bool
	allowedTypes      map[CodegenLogType]bool
	preface           bool
}

type CodegenLogCategory string
type CodegenLogType string

const (
	// CLogLoaders Logger category for loaders. All loader logging should use this category
	CLogLoaders CodegenLogCategory = "loaders"

	// CLogGenerators Logger category for generators. All generator logging should use this category
	CLogGenerators CodegenLogCategory = "generators"

	// CLogDeployers Logger category for deployers. All deployer logging should use this category
	CLogDeployers CodegenLogCategory = "deployers"

	// CLogCore Logger category for Core codegen. All core logging should use this category
	CLogCore CodegenLogCategory = "core"

	CLogLog   CodegenLogType = "log"
	CLogError CodegenLogType = "error"
	CLogFatal CodegenLogType = "fatal"
)

// CLog an exported single instance variable usable anywhere. Will be insantiated in the init() func
var CLog *CodegenLogger

func init() {
	CLog = NewLogger()
	// Add default categories. Any code can remove these and/or add other categories of their own as deemed necessary
	CLog.AddCategory(CLogLoaders)
	CLog.AddCategory(CLogGenerators)
	CLog.AddCategory(CLogDeployers)
	CLog.AddCategory(CLogCore)
}

// NewLogger
// While used in the init() func here to create a single usable exported logger.. this method can be used
// to create more instances of loggers if needed. It creates the empty Map of logger category types and
// sets the preface to true.
func NewLogger() *CodegenLogger {
	return &CodegenLogger{
		categories:        make([]CodegenLogCategory, 0),
		allowedCategories: make(map[CodegenLogCategory]bool),
		allowedTypes:      make(map[CodegenLogType]bool),
		preface:           true,
	}
}

func (l *CodegenLogger) AddCategory(category CodegenLogCategory) {
	l.categories = append(l.categories, category)
}

func (l *CodegenLogger) GetCategory(category string) *CodegenLogCategory {
	for _, cat := range l.categories {
		if category == string(cat) {
			return &cat
		}
	}

	return nil
}

// AllowType
// This function takes in a string value and if it finds a matching string log type to the value provided will add it
// to the list of allowed types to be logged.
func (l *CodegenLogger) AllowType(typ string) {
	switch typ {
	case "log":
		l.allowedTypes[CLogLog] = true
	case "error":
		l.allowedTypes[CLogError] = true
	case "fatal":
		l.allowedTypes[CLogFatal] = true
	}
}

// AllowCategory
// Adds a new category to the allowed logging map. The category must be of type CodegenLogType (a string)
// Additional categories can be added outside of the default "core" categories defined above
func (l *CodegenLogger) AllowCategory(category CodegenLogCategory) {
	for _, cat := range l.categories {
		if cat == category {
			l.allowedCategories[category] = true
			return
		}
	}
}

// RemoveAllowedCategory
// Removes a category from the map of logging categories. The removal of a category prevents any calls to
// Log or Fatal methods with the category removed from being displayed to the console. This is how to hide
// logging output as needed.
func (l *CodegenLogger) RemoveAllowedCategory(category CodegenLogCategory) {
	if _, ok := l.allowedCategories[category]; ok {
		delete(l.allowedCategories, category)
	}
}

// RemoveAllAllowed
// This will remove all current allowed categories.. good to clear them out before adding specific categories
// for logging output.
func (l *CodegenLogger) RemoveAllAllowed() {
	for lg := range l.allowedCategories {
		delete(l.allowedCategories, lg)
	}
}

// Log
// This call simply redirects the format and v parameters to the standard library log.Printf() call. If the
// logger preface is true, it will preface the output with Codegen Logger: to help differentiate this logging
// from other logging output that may be used elsewhere. It will also include the preface of category
func (l *CodegenLogger) Log(category CodegenLogCategory, source string, format string, v ...any) {
	if _, ok := l.allowedTypes[CLogLog]; ok {
		if _, ok = l.allowedCategories[category]; ok {
			if l.preface {
				log.Printf("Codegen Logger: [%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v))
			} else {
				log.Printf("[%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v))
			}
		}
	}
}

// Fatal
// This call simply redirects the format and v parameters to the standard library log.Fatalf() call. If the
// logger preface is true, it will preface the output with Codegen Logger: to help differentiate this logging
// from other logging output that may be used elsewhere. It will also include the preface of category
func (l *CodegenLogger) Fatal(category CodegenLogCategory, source string, format string, v ...any) {
	if _, ok := l.allowedTypes[CLogFatal]; ok {
		if _, ok = l.allowedCategories[category]; ok {
			if l.preface {
				log.Fatalf("Codegen Logger: [%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v))
			} else {
				log.Fatalf("[%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v))
			}
		}
	}
}

// Fatal
// This call simply redirects the format and v parameters to the standard library log.Fatalf() call. If the
// logger preface is true, it will preface the output with Codegen Logger: to help differentiate this logging
// from other logging output that may be used elsewhere. It will also include the preface of category
func (l *CodegenLogger) Error(category CodegenLogCategory, source string, format string, v ...any) {
	if _, ok := l.allowedTypes[CLogError]; ok {
		if _, ok = l.allowedCategories[category]; ok {
			if l.preface {
				log.Printf("Codegen Logger: ERROR: [%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v...))
			} else {
				log.Printf("ERROR: [%s]-%s %s", strings.ToUpper(string(category)), source, fmt.Sprintf(format, v...))
			}
		}
	}
}
