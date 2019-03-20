// +build debug

package logging

func (logger *concreteLogger) Debug(v ...interface{}) {
	logger.Debugd(3, v...)
}

func (logger *concreteLogger) Debugf(format string, v ...interface{}) {
	logger.Debugdf(3, format, v...)
}
