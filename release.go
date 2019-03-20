// +build !debug

package logging

func (logger *concreteLogger) Debug(v ...interface{}) {
}

func (logger *concreteLogger) Debugf(format string, v ...interface{}) {
}
