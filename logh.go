package logh

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sort"
	"sync"

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

	mutex sync.Mutex
)

type keyValue struct {
	key    string
	value  interface{}
	rvalue reflect.Value
	kind   reflect.Kind
}

func newItem(key string, value interface{}) keyValue {

	rvalue := reflect.ValueOf(value)

	return keyValue{
		key:    key,
		value:  value,
		rvalue: rvalue,
		kind:   rvalue.Kind(),
	}
}

// byKey - implements sort.Interface based on the Key field.
type byKey []keyValue

// Len - returns the length
func (s byKey) Len() int {

	return len(s)
}

// Less - compares
func (s byKey) Less(i, j int) bool {

	return (len(s[i].key) == len(s[j].key)) && s[i].key < s[j].key
}

// Swap - swap indexes
func (s byKey) Swap(i, j int) {

	s[i], s[j] = s[j], s[i]
}

// ContextualLogger - a struct containing all valid event loggers (each one can be null if not enabled)
type ContextualLogger struct {
	keyValues []keyValue
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

	return cl.ErrorLineC(2)
}

// ErrorLineC - returns the event logger using the configured context
func (cl *ContextualLogger) ErrorLineC(skippedStackFrames int) *zerolog.Event {

	_, filename, line, ok := runtime.Caller(skippedStackFrames)
	ev := Error()
	if !ok {
		filename = "unknown"
		line = -1
	}

	ev = ev.Str("@file", filename)

	if ok {
		ev = ev.Int("@line", line)
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

	mutex.Lock()
	defer mutex.Unlock()

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

func toKeyValueArray(keyValues ...interface{}) []keyValue {

	numKeyValues := len(keyValues)

	items := make([]keyValue, numKeyValues/2)
	j := 0
	for i := 0; i < numKeyValues; i += 2 {
		items[j] = newItem(keyValues[i].(string), keyValues[i+1])
		j++
	}

	return items
}

// CreateContextualLogger - creates loggers with context
func CreateContextualLogger(keyValues ...interface{}) *ContextualLogger {

	numKeyValues := len(keyValues)
	if numKeyValues%2 != 0 {
		panic(ErrWrongNumberOfArgs)
	}

	kvs := toKeyValueArray(keyValues...)
	sort.Sort(byKey(kvs))

	return &ContextualLogger{
		keyValues: kvs,
	}
}

// Append - appends more context
func (cl *ContextualLogger) Append(keyValues ...interface{}) error {

	numKeyValues := len(keyValues)
	if numKeyValues%2 != 0 {
		return ErrWrongNumberOfArgs
	}

	cl.keyValues = append(cl.keyValues, toKeyValueArray(keyValues...)...)

	sort.Sort(byKey(cl.keyValues))

	return nil
}

// MustAppend - appends more context, panics if any error is founbd
func (cl *ContextualLogger) MustAppend(keyValues ...interface{}) {

	err := cl.Append(keyValues...)
	if err != nil {
		panic(err)
	}
}

// addContext - add event logger context
func (cl *ContextualLogger) addContext(eventlLogger *zerolog.Event) *zerolog.Event {

	if eventlLogger == nil {
		return nil
	}

	for j := 0; j < len(cl.keyValues); j++ {

		switch cl.keyValues[j].kind {

		case reflect.String:

			eventlLogger = eventlLogger.Str(cl.keyValues[j].key, cl.keyValues[j].rvalue.String())

		case reflect.Int:

			eventlLogger = eventlLogger.Int(cl.keyValues[j].key, int(cl.keyValues[j].rvalue.Int()))

		case reflect.Int8:

			eventlLogger = eventlLogger.Int8(cl.keyValues[j].key, int8(cl.keyValues[j].rvalue.Int()))

		case reflect.Int16:

			eventlLogger = eventlLogger.Int16(cl.keyValues[j].key, int16(cl.keyValues[j].rvalue.Int()))

		case reflect.Int32:

			eventlLogger = eventlLogger.Int32(cl.keyValues[j].key, int32(cl.keyValues[j].rvalue.Int()))

		case reflect.Int64:

			eventlLogger = eventlLogger.Int64(cl.keyValues[j].key, cl.keyValues[j].rvalue.Int())

		case reflect.Uint:

			eventlLogger = eventlLogger.Uint(cl.keyValues[j].key, uint(cl.keyValues[j].rvalue.Uint()))

		case reflect.Uint8:

			eventlLogger = eventlLogger.Uint8(cl.keyValues[j].key, uint8(cl.keyValues[j].rvalue.Uint()))

		case reflect.Uint16:

			eventlLogger = eventlLogger.Uint16(cl.keyValues[j].key, uint16(cl.keyValues[j].rvalue.Uint()))

		case reflect.Uint32:

			eventlLogger = eventlLogger.Uint32(cl.keyValues[j].key, uint32(cl.keyValues[j].rvalue.Uint()))

		case reflect.Uint64:

			eventlLogger = eventlLogger.Uint64(cl.keyValues[j].key, cl.keyValues[j].rvalue.Uint())

		case reflect.Float32:

			eventlLogger = eventlLogger.Float32(cl.keyValues[j].key, float32(cl.keyValues[j].rvalue.Float()))

		case reflect.Float64:

			eventlLogger = eventlLogger.Float64(cl.keyValues[j].key, cl.keyValues[j].rvalue.Float())

		case reflect.Bool:

			eventlLogger = eventlLogger.Bool(cl.keyValues[j].key, cl.keyValues[j].rvalue.Bool())

		default:

			eventlLogger = eventlLogger.Interface(cl.keyValues[j].key, cl.keyValues[j].rvalue.Interface())
		}
	}

	return eventlLogger
}

// GetContexts - returns the logger contexts
func (cl *ContextualLogger) GetContexts() []interface{} {

	oldKVs := make([]interface{}, len(cl.keyValues)*2)

	i := 0
	for _, item := range cl.keyValues {
		oldKVs[i] = item.key
		oldKVs[i+1] = item.value
		i += 2
	}

	return oldKVs
}

// CreateFromContext - creates a new logger context from this context
func (cl *ContextualLogger) CreateFromContext(keyValues ...interface{}) (*ContextualLogger, error) {

	oldKVs := cl.GetContexts()

	ccl := CreateContextualLogger(oldKVs...)
	err := ccl.Append(keyValues...)
	if err != nil {
		return nil, err
	}

	return ccl, nil
}

// MustCreateFromContext - creates a new logger context from this context, raises panic if some error
func (cl *ContextualLogger) MustCreateFromContext(keyValues ...interface{}) *ContextualLogger {

	ccl, err := cl.CreateFromContext(keyValues...)
	if err != nil {
		panic(err)
	}

	return ccl
}
