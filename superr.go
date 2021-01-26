package superr

import (
	"github.com/jimlawless/whereami"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type Op string
type Kind string
type Message string

// Application error codes
const (
	ECONFLICT Kind = "conflict"  // action cannot be performed
	EINTERNAL      = "internal"  // internal error
	EINVALID       = "invalid"   // validation failed
	ENOTFOUND      = "not_found" // entity does not exist
)

type Error struct {
	// Machine-readable error code.
	Kind Kind // code of error: like conflict, invaild, not_found etc

	// Human-readable message.
	Message Message

	// Logical operation and nested error.
	Op  Op // Operation is a unique string describing a method or a function
	Err error

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

	entry := log.WithFields(log.Fields{
		"stackTrace": Ops(superError),
	})

	switch superError.Severity {
	default:
		entry.Info(superError.Message)
	}
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

func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Kind != "" {
		return string(e.Kind)
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return EINTERNAL
}

func (e *Error) Error() string {
	return string(e.Message)
}
