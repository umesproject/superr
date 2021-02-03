package superr

import (
	"github.com/jimlawless/whereami"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// Op stands for Operation. Is a unique string describing a method or a function
type Op string

// Kind groups all errors into smaller categories
// Can be predefined codes ( like http / gRPC )
type Kind uint32

// Message is a Human-readable message.
type Message string

// Field is used for passing extra data to the error ( like variables content usefull for debug  and query on log )
type Fields map[string]interface{}

type Error struct {
	// Machine-readable error code.
	Kind Kind // code of error: like conflict, invaild, not_found etc

	Message Message

	// Logical operation and nested error.
	Op  Op
	Err error

	ExtraFields Fields // Extra fields useful for logging

	// Metadata
	Severity log.Level
	Caller   string
}

func E(args ...interface{}) error {
	e := &Error{}

	e.Caller = whereami.WhereAmI(2)

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Message:
			e.Message = arg
		case error:
			e.Err = arg
		case Kind:
			e.Kind = arg
		case Fields:
			e.ExtraFields = arg
		default:
			panic("bad call to serros.")
		}
	}
	return e
}

func Log(e error) {
	superError, ok := e.(*Error)
	if !ok {
		log.Error(e)
		return
	}

	fields := log.Fields{
		"stackTrace": Ops(superError),
		"caller":     Caller(superError),
	}

	// Get fields from first error
	/*	for fieldName, fieldValue := range getFields(superError) {
		fields[fieldName] = fieldValue
	}*/

	fields["extraFields"] = getFields(superError)

	entry := log.WithFields(fields)

	switch superError.Severity {
	default:
		entry.Info(superError.Message)
	}
}

// Fields return the extra fields added to the error
func getFields(e *Error) Fields {
	subErr, ok := e.Err.(*Error)
	if !ok {
		return e.ExtraFields
	}

	return getFields(subErr)
}

// Ops returns the "stack" of operations
// for each generated error
func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)
	return res
}

// ErrorCode returns the error code associated with the error
func ErrorCode(err error) Kind {
	e, ok := err.(*Error)
	if !ok {
		return 0
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return ErrorCode(e.Err)
}

func Caller(e *Error) string {
	return e.Caller
}

func Callers(e *Error) []string {
	res := []string{}
	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Callers(subErr)...)
	return res
}

func (e *Error) Error() string {
	return string(e.Message)
}
