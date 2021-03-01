package superr

import (
	"fmt"
	"github.com/jimlawless/whereami"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Op stands for Operation. Is a unique string describing a method or a function
type Op string

// Kind groups all errors into smaller categories
// Can be predefined codes ( like http / gRPC )
type Kind uint32

// Message is a Human-readable message.
type Message string

// Field is used for passing extra data to the error ( like variables content usefull for debug  and query on log )
type Fields map[string]interface{}

var istance *zap.Logger

type Severity int8
type logLevel int8

const (
	InfoLevel logLevel = iota
	DebugLevel
	ErrorLevel
)

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityError
)

func init() {
	Init(InfoLevel)
}

func Init(level logLevel) {

	choosedLevel := zap.InfoLevel

	switch level {
	case DebugLevel:
		choosedLevel = zap.DebugLevel
		break
	case ErrorLevel:
		choosedLevel = zap.ErrorLevel
		break
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(choosedLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "severity",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	ist, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	istance = ist
}

type Error struct {
	// Machine-readable error code.
	Kind Kind // code of error: like conflict, invaild, not_found etc

	Message Message

	// Logical operation and nested error.
	Op  Op
	Err error

	ExtraFields Fields // Extra fields useful for logging

	// Metadata
	Severity Severity
	Caller   string
}

func E(args ...interface{}) error {
	e := &Error{}

	e.Caller = whereami.WhereAmI(2)
	e.Severity = SeverityInfo

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
		case Severity:
			e.Severity = arg
		default:
			panic(fmt.Sprintf("bad call to superr.E: not recognised type %v", arg))
		}
	}
	return e
}

func ErrorWithLog(args ...interface{}) error {
	e := E(args...)
	Log(e)
	return e
}

func Log(e error) {
	superError, ok := e.(*Error)
	if !ok {
		log.Error(e)
		return
	}

	fields := []zap.Field{
		zap.Any("stackTrace", Ops(superError)),
		zap.Any("caller", Caller(superError))}

	// Get fields from first error
	/*	for fieldName, fieldValue := range getFields(superError) {
		fields[fieldName] = fieldValue
	}*/

	fields = append(fields, zap.Any("extraFields", getFields(superError)))

	entry := istance.With(fields...).WithOptions(zap.WithCaller(false))

	errorMessage := string(superError.Message)
	switch superError.Severity {
	case SeverityDebug:
		entry.Debug(errorMessage)
		break
	case SeverityError:
		entry.Error(errorMessage)
	default:
		entry.Info(string(superError.Message))
		break
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
