package logh

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"

	"github.com/rs/zerolog"
)

//
// Has some useful logging functions.
// logh -> log helper
// @author rnojiri
//

// Level - type
type Level string

const (
	// INFO - log level
	INFO Level = "info"

	// DEBUG - log level
	DEBUG Level = "debug"

	// WARN - log level
	WARN Level = "warn"

	// ERROR - log level
	ERROR Level = "error"

	// FATAL - log level
	FATAL Level = "fatal"

	// PANIC - log level
	PANIC Level = "panic"

	// NONE - log level
	NONE Level = "none"

	// SILENT - log level
	SILENT Level = "silent"
)

// Format - the logger's output format
type Format string

const (
	// JSON - json format
	JSON Format = "json"

	// CONSOLE - plain text format
	CONSOLE Format = "console"
)

var (
	logger zerolog.Logger

	// InfoEnabled - check if this level is enabled
	InfoEnabled bool

	// DebugEnabled - check if this level is enabled
	DebugEnabled bool

	// WarnEnabled - check if this level is enabled
	WarnEnabled bool

	// ErrorEnabled - check if this level is enabled
	ErrorEnabled bool

	// FatalEnabled - check if this level is enabled
	FatalEnabled bool

	// PanicEnabled - check if this level is enabled
	PanicEnabled bool

	// ErrWrongNumberOfArgs ...
	ErrWrongNumberOfArgs error = errors.New("the number of arguments must be even")
)

// ContextualLogger - a struct containing all valid event loggers (each one can be null if not enabled)
type ContextualLogger struct {
	numKeyValues int
	keyValues    []interface{}
}

// Info - returns the event logger using the configured context
func (cl *ContextualLogger) Info() *zerolog.Event {
	return cl.addContext(Info())
}

// Debug - returns the event logger using the configured context
func (cl *ContextualLogger) Debug() *zerolog.Event {
	return cl.addContext(Debug())
}

// Warn - returns the event logger using the configured context
func (cl *ContextualLogger) Warn() *zerolog.Event {
	return cl.addContext(Warn())
}

// Error - returns the event logger using the configured context
func (cl *ContextualLogger) Error() *zerolog.Event {
	return cl.addContext(Error())
}

// ErrorLine - returns the event logger using the configured context
func (cl *ContextualLogger) ErrorLine() *zerolog.Event {

	_, filename, line, ok := runtime.Caller(1)
	ev := Error()
	if !ok {
		filename = "unknown"
		line = -1
	}

	ev = ev.Str("_file_", filename)

	if ok {
		ev = ev.Int("_line_", line)
	}

	return cl.addContext(ev)
}

// Fatal - returns the event logger using the configured context
func (cl *ContextualLogger) Fatal() *zerolog.Event {
	return cl.addContext(Fatal())
}

// Panic - returns the event logger using the configured context
func (cl *ContextualLogger) Panic() *zerolog.Event {
	return cl.addContext(Panic())
}

// ConfigureGlobalLogger - configures the logger globally
func ConfigureGlobalLogger(lvl Level, fmt Format) *zerolog.Logger {

	return ConfigureCustomLogger(lvl, fmt, os.Stdout)
}

// ConfigureCustomLogger - configures the logger globally
func ConfigureCustomLogger(lvl Level, fmt Format, out io.Writer) *zerolog.Logger {

	switch lvl {
	case INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case PANIC:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case FATAL:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case NONE:
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case SILENT:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	var writer io.Writer

	if fmt == CONSOLE {
		writer = zerolog.ConsoleWriter{Out: out}
	} else {
		writer = out
	}

	logger = zerolog.New(writer).With().Timestamp().Logger()

	InfoEnabled = Info() != nil
	DebugEnabled = Debug() != nil
	WarnEnabled = Warn() != nil
	ErrorEnabled = Error() != nil
	PanicEnabled = Panic() != nil
	FatalEnabled = Fatal() != nil

	return &logger
}

// SendToStdout - logs a output with no log format
func SendToStdout(output string) {

	fmt.Println(output)
}

// Info - returns the info event logger if any
func Info() *zerolog.Event {
	if e := logger.Info(); e.Enabled() {
		return e
	}
	return nil
}

// Debug - returns the debug event logger if any
func Debug() *zerolog.Event {
	if e := logger.Debug(); e.Enabled() {
		return e
	}
	return nil
}

// Warn - returns the error event logger if any
func Warn() *zerolog.Event {
	if e := logger.Warn(); e.Enabled() {
		return e
	}
	return nil
}

// Error - returns the error event logger if any
func Error() *zerolog.Event {
	if e := logger.Error(); e.Enabled() {
		return e
	}
	return nil
}

// Panic - returns the error event logger if any
func Panic() *zerolog.Event {
	if e := logger.Panic(); e.Enabled() {
		return e
	}
	return nil
}

// Fatal - returns the error event logger if any
func Fatal() *zerolog.Event {
	if e := logger.Fatal(); e.Enabled() {
		return e
	}
	return nil
}

// Logger - returns the logger itself
func Logger() *zerolog.Logger {

	return &logger
}

// CreateContextualLogger - creates loggers with context
func CreateContextualLogger(keyValues ...interface{}) *ContextualLogger {

	numKeyValues := len(keyValues)
	if numKeyValues%2 != 0 {
		panic(ErrWrongNumberOfArgs)
	}

	return &ContextualLogger{
		numKeyValues: numKeyValues,
		keyValues:    keyValues,
	}
}

// Append - appends more context
func (cl *ContextualLogger) Append(keyValues ...interface{}) error {

	numKeyValues := len(keyValues)
	if numKeyValues%2 != 0 {
		return ErrWrongNumberOfArgs
	}

	cl.keyValues = append(cl.keyValues, keyValues...)
	cl.numKeyValues += numKeyValues

	return nil
}

// addContext - add event logger context
func (cl *ContextualLogger) addContext(eventlLogger *zerolog.Event) *zerolog.Event {

	if eventlLogger == nil {
		return nil
	}

	for j := 0; j < cl.numKeyValues; j += 2 {

		key := cl.keyValues[j].(string)
		value := reflect.ValueOf(cl.keyValues[j+1])

		switch value.Kind() {

		case reflect.String:

			eventlLogger = eventlLogger.Str(key, value.String())

		case reflect.Int:

			eventlLogger = eventlLogger.Int(key, int(value.Int()))

		case reflect.Int8:

			eventlLogger = eventlLogger.Int8(key, int8(value.Int()))

		case reflect.Int16:

			eventlLogger = eventlLogger.Int16(key, int16(value.Int()))

		case reflect.Int32:

			eventlLogger = eventlLogger.Int32(key, int32(value.Int()))

		case reflect.Int64:

			eventlLogger = eventlLogger.Int64(key, value.Int())

		case reflect.Uint:

			eventlLogger = eventlLogger.Uint(key, uint(value.Uint()))

		case reflect.Uint8:

			eventlLogger = eventlLogger.Uint8(key, uint8(value.Uint()))

		case reflect.Uint16:

			eventlLogger = eventlLogger.Uint16(key, uint16(value.Uint()))

		case reflect.Uint32:

			eventlLogger = eventlLogger.Uint32(key, uint32(value.Uint()))

		case reflect.Uint64:

			eventlLogger = eventlLogger.Uint64(key, value.Uint())

		case reflect.Float32:

			eventlLogger = eventlLogger.Float32(key, float32(value.Float()))

		case reflect.Float64:

			eventlLogger = eventlLogger.Float64(key, value.Float())

		case reflect.Bool:

			eventlLogger = eventlLogger.Bool(key, value.Bool())

		default:

			eventlLogger = eventlLogger.Interface(key, value.Interface())
		}
	}

	return eventlLogger
}

// GetContexts - returns the logger contexts
func (cl *ContextualLogger) GetContexts() []interface{} {

	return cl.keyValues
}

// CreateFromContext - creates a new logger context from this context
func (el *ContextualLogger) CreateFromContext(keyValues ...interface{}) (*ContextualLogger, error) {

	cl := CreateContextualLogger(el.keyValues...)
	err := cl.Append(keyValues...)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// MustCreateFromContext - creates a new logger context from this context, raises panic if some error
func (el *ContextualLogger) MustCreateFromContext(keyValues ...interface{}) *ContextualLogger {

	cl, err := el.CreateFromContext(keyValues...)
	if err != nil {
		panic(err)
	}

	return cl
}
