package logging

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DefaultLoggerSuite struct {
	suite.Suite
}

func (suite *DefaultLoggerSuite) TestBufferFatal() {
	logger := DefaultLogger()

	var buffer bytes.Buffer
	concrete, ok := defaultLogger.value.(*concreteLogger)
	suite.Require().True(ok)
	concrete.SetOutput(&buffer)

	msg0 := "dfdffdc4-7281-4b77-8ed7-9b07a10ab354"
	pc, file, line, ok := runtime.Caller(0)
	suite.Require().True(ok)
	logger.Fatal(msg0)

	suite.Require().Contains(buffer.String(), "Fatal")
	suite.Require().Contains(buffer.String(), msg0)
	funcName := runtime.FuncForPC(pc).Name()
	file = filepath.Base(file)
	suite.Require().Contains(buffer.String(), funcName)
	suite.Require().Contains(buffer.String(), file)
	suite.Require().Contains(buffer.String(), fmt.Sprintf("%d", line+2))
}

func TestDefaultLogger(t *testing.T) {
	suite.Run(t, new(DefaultLoggerSuite))
}
