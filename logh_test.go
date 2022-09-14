package logh_test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/uol/logh"
)

type testSuite struct {
	suite.Suite
	loggerBuffer *logh.StringWriter
}

func (ts *testSuite) SetupTest() {

	ts.loggerBuffer = logh.NewStringWriter(256)
	logh.ConfigureCustomLogger(logh.INFO, logh.JSON, ts.loggerBuffer)
}

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

func (ts *testSuite) baseLoggerTest(contexts []interface{}) {

	cl := logh.CreateContextualLogger(contexts...)

	cl.Info().Msg("hello world")

	storedContexts := cl.GetContexts()
	ts.Equal(contexts, storedContexts, "expected same contexts")
}

func (ts *testSuite) TestContextualLogger() {

	logh.ConfigureGlobalLogger(logh.INFO, logh.CONSOLE)

	contexts := []interface{}{"context1", "test1", "context2", 2}

	ts.baseLoggerTest(contexts)
}

func (ts *testSuite) testBufferContents(expected string) {

	ts.Equal(expected, strings.Trim(string(ts.loggerBuffer.Bytes()), "\n"), "expected same log message")
}

func (ts *testSuite) TestContextualCustomLogger() {

	contexts := []interface{}{"context3", "test3", "context4", true}

	now := time.Now()

	ts.baseLoggerTest(contexts)

	expected := fmt.Sprintf(`{"level":"info","context3":"test3","context4":true,"time":"%s","message":"hello world"}`, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func (ts *testSuite) TestContextualLoggerAppend() {

	contexts := []interface{}{"context5", "test5"}

	now := time.Now()

	cl := logh.CreateContextualLogger(contexts...)

	err := cl.Append("context6", 6, "context7", 0.7)
	ts.NoError(err, "expects no errors")

	cl.Info().Msg("append test")

	expected := fmt.Sprintf(`{"level":"info","context5":"test5","context6":6,"context7":0.7,"time":"%s","message":"append test"}`, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func (ts *testSuite) TestContextualLoggerMustAppend() {

	contexts := []interface{}{"context8", "test8"}

	now := time.Now()

	cl := logh.CreateContextualLogger(contexts...)

	cl.MustAppend("context9", 9, "context10", 0.1)

	cl.Info().Msg("must append test")

	expected := fmt.Sprintf(`{"level":"info","context8":"test8","context9":9,"context10":0.1,"time":"%s","message":"must append test"}`, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func (ts *testSuite) TestCreateFromContext() {

	cl := logh.CreateContextualLogger("root_key", "root_val")

	now := time.Now()

	ncl, err := cl.CreateFromContext("node_key1", 1, "node_key2", 2)
	ts.NoError(err, "expects no errors")

	ncl.Info().Msg("create from context")

	expected := fmt.Sprintf(`{"level":"info","root_key":"root_val","node_key1":1,"node_key2":2,"time":"%s","message":"create from context"}`, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func (ts *testSuite) TestMustCreateFromContext() {

	cl := logh.CreateContextualLogger("root_key", "root_val")

	now := time.Now()

	ncl := cl.MustCreateFromContext("node_key3", 3, "node_key4", 4)

	ncl.Info().Msg("must create from context")

	expected := fmt.Sprintf(`{"level":"info","root_key":"root_val","node_key3":3,"node_key4":4,"time":"%s","message":"must create from context"}`, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func getFileAndLine() (string, int) {

	_, filename, line, _ := runtime.Caller(1)

	return filename, line
}

func (ts *testSuite) TestErrorLine() {

	now := time.Now()

	logger := logh.CreateContextualLogger("pkg", "logh_test")

	expectedFilename, expectedLine := getFileAndLine()
	expectedLine += 2
	logger.ErrorLine().Err(errors.New("test error")).Msg("message")

	expected := fmt.Sprintf(`{"level":"error","@file":"%s","@line":%d,"pkg":"logh_test","error":"test error","time":"%s","message":"message"}`, expectedFilename, expectedLine, now.Format(time.RFC3339))

	ts.testBufferContents(expected)
}

func TestSuite(t *testing.T) {

	suite.Run(t, new(testSuite))
}
