package httpway

import "fmt"

//logger interface
type Logger interface {
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
}

const default_depth = 5

type internalLogger struct {
	l      Logger
	id     uint64
	prefix string
}

func (il *internalLogger) Info(v ...interface{}) {
	if il.l == nil {
		return
	}
	il.log(il.l.Info, v...)
}

func (il *internalLogger) Warning(v ...interface{}) {
	if il.l == nil {
		return
	}

	il.log(il.l.Warning, v...)
}

func (il *internalLogger) Error(v ...interface{}) {
	if il.l == nil {
		return
	}

	il.log(il.l.Error, v...)
}

func (il *internalLogger) Debug(v ...interface{}) {
	if il.l == nil {
		return
	}
	il.log(il.l.Debug, v...)
}

func (il *internalLogger) log(f func(v ...interface{}), v ...interface{}) {
	if len(v) > 1 {
		v[0] = fmt.Sprintf("[%x] %s", il.id, v[0].(string))
	} else {
		v[0] = fmt.Sprintf("[%x] %v", il.id, v[0])
	}

	if il.prefix != "" {
		v[0] = fmt.Sprintf("[%s] %s", il.prefix, v[0].(string))
	}

	f(v...)
}

type internalServerLoggerWriter struct {
	l Logger
}

func (islw *internalServerLoggerWriter) Write(b []byte) (n int, err error) {
	islw.l.Error("Go Server: %s", string(b))
	return len(b), nil
}
