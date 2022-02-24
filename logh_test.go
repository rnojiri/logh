package logh_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uol/logh"
)

// TestGlobalConfiguration - tests the global configuration
func TestGlobalConfiguration(t *testing.T) {

	logh.ConfigureGlobalLogger(logh.INFO, logh.CONSOLE)

	assert.True(t, logh.Info() != nil, "expected true")
	assert.False(t, logh.Debug() != nil, "expected false")
	assert.True(t, logh.Warn() != nil, "expected true")
	assert.True(t, logh.Error() != nil, "expected true")
	assert.True(t, logh.Fatal() != nil, "expected true")
	assert.True(t, logh.Panic() != nil, "expected true")
	assert.True(t, logh.Logger() != nil, "expected true")
}

func baseLoggerTest(t *testing.T, contexts []interface{}) {

	cl := logh.CreateContextualLogger(contexts...)

	cl.Info().Msg("hello world")

	storedContexts := cl.GetContexts()
	assert.Equal(t, contexts, storedContexts, "expected same contexts")
}

func TestContextualLogger(t *testing.T) {

	logh.ConfigureGlobalLogger(logh.INFO, logh.CONSOLE)

	contexts := []interface{}{"context1", "test1", "context2", 2}

	baseLoggerTest(t, contexts)
}

func TestContextualCustomLogger(t *testing.T) {

	writer := logh.NewStringWriter(256)

	logh.ConfigureCustomLogger(logh.INFO, logh.JSON, writer)

	contexts := []interface{}{"context3", "test3", "context4", true}

	now := time.Now()

	baseLoggerTest(t, contexts)

	expected := fmt.Sprintf(`{"level":"info","context3":"test3","context4":true,"time":"%s","message":"hello world"}`, now.Format(time.RFC3339))

	assert.Equal(t, expected, strings.Trim(string(writer.Bytes()), "\n"), "expected same log message")
}

func TestContextualLoggerAppend(t *testing.T) {

	writer := logh.NewStringWriter(256)

	logh.ConfigureCustomLogger(logh.INFO, logh.JSON, writer)

	contexts := []interface{}{"context5", "test5"}

	now := time.Now()

	cl := logh.CreateContextualLogger(contexts...)

	err := cl.Append("context6", 6, "context7", 0.7)
	assert.NoError(t, err, "expects no errors")

	cl.Info().Msg("append test")

	expected := fmt.Sprintf(`{"level":"info","context5":"test5","context6":6,"context7":0.7,"time":"%s","message":"append test"}`, now.Format(time.RFC3339))

	assert.Equal(t, expected, strings.Trim(string(writer.Bytes()), "\n"), "expected same log message")
}
