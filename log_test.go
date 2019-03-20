package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/suite"
)

type DefaultLoggerSuite struct {
	suite.Suite
}

func (suite *DefaultLoggerSuite) TestFatal() {
	color.NoColor = true
	defer func() {
		color.NoColor = false
	}()
	logger := DefaultLogger()

	var buffer strings.Builder
	concrete, ok := logger.(*concreteLogger)
	suite.Require().True(ok)
	concrete.SetOutput(&buffer)
	pc, file, line, ok := runtime.Caller(0)
	suite.Require().True(ok)
	logger.Fatal()
	output := buffer.String()
	suite.Require().Contains(output, level.fatal.name)
	funcName := runtime.FuncForPC(pc).Name()
	file = filepath.Base(file)
	suite.Require().Contains(output, funcName)
	suite.Require().Contains(output, file)
	suite.Require().Contains(output, fmt.Sprintf("%d", line+2))

	msg0 := "dfdffdc4-7281-4b77-8ed7-9b07a10ab354"
	buffer.Reset()
	pc, file, line, ok = runtime.Caller(0)
	suite.Require().True(ok)
	logger.Fatal(msg0)
	output = buffer.String()
	suite.Require().Contains(output, level.fatal.name)
	suite.Require().Contains(output, msg0)
	funcName = runtime.FuncForPC(pc).Name()
	file = filepath.Base(file)
	suite.Require().Contains(output, funcName)
	suite.Require().Contains(output, file)
	suite.Require().Contains(output, fmt.Sprintf("%d", line+2))
}

func (suite *DefaultLoggerSuite) TestErrorField() {
	color.NoColor = true
	defer func() {
		color.NoColor = false
	}()
	logger := DefaultLogger()

	var buffer strings.Builder
	concrete, ok := logger.(*concreteLogger)
	suite.Require().True(ok)
	concrete.SetOutput(&buffer)
	logger.Error()
	output := buffer.String()
	suite.Require().Contains(output, level.error.name)

	msg0 := "dfdffdc4-7281-4b77-8ed7-9b07a10ab354"
	buffer.Reset()
	logger.Error(msg0)
	output = buffer.String()
	suite.Require().Contains(output, level.error.name)
	suite.Require().Contains(output, msg0)

	buffer.Reset()
	logger.Error(Field{})
	output = buffer.String()
	suite.Require().Contains(output, level.error.name)

	field0 := Field{
		"bool": true,
		"int":  10,
	}
	buffer.Reset()
	logger.Error(field0)
	output = buffer.String()
	suite.Require().Contains(output, level.error.name)
	for key, value := range field0 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}

	msg1 := "2c9a3582-2990-42c5-89cf-4c6f9cde4e1e"
	field1 := Field{
		"float":  23.4,
		"string": "abc",
	}
	buffer.Reset()
	logger.Error(msg0, field0, msg1, &field1)
	output = buffer.String()
	suite.Require().Contains(output, level.error.name)
	suite.Require().Contains(output, msg0)
	suite.Require().Contains(output, msg1)
	for key, value := range field0 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}
	for key, value := range field1 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}
}

func (suite *DefaultLoggerSuite) TestWarnfField() {
	color.NoColor = true
	defer func() {
		color.NoColor = false
	}()
	logger := DefaultLogger()

	var buffer strings.Builder
	concrete, ok := logger.(*concreteLogger)
	suite.Require().True(ok)
	concrete.SetOutput(&buffer)
	msg0 := "dfdffdc4-7281-4b77-8ed7-9b07a10ab354"
	logger.Errorf(msg0)
	output := buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)

	buffer.Reset()
	logger.Errorf("%s", msg0)
	output = buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)

	logger.Errorf(msg0, Field{})
	output = buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)

	field0 := Field{
		"bool": true,
		"int":  10,
	}
	buffer.Reset()
	logger.Errorf(msg0, field0)
	output = buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)
	for key, value := range field0 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}

	buffer.Reset()
	logger.Errorf("%s", msg0, field0)
	output = buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)
	for key, value := range field0 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}

	msg1 := "2c9a3582-2990-42c5-8cf-4c6f9cde4e1e"
	msg2 := "5ac9795c-2c06-43b1-aae8-d72a8f573738"
	field1 := Field{
		"float":  23.4,
		"string": "abc",
	}
	buffer.Reset()
	logger.Errorf("%s, %s", msg0, msg1, field0, msg2, &field1)
	output = buffer.String()
	suite.Require().Contains(output, level.warn.name)
	suite.Require().Contains(output, msg0)
	suite.Require().Contains(output, msg1)
	suite.Require().Contains(output, msg2)
	for key, value := range field0 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}
	for key, value := range field1 {
		suite.Require().Contains(output, fmt.Sprintf("%s(%v)", key, value))
	}
}

func TestDefaultLogger(t *testing.T) {
	suite.Run(t, new(DefaultLoggerSuite))
}
