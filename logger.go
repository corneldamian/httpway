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

type internalLogger struct{
	l Logger
	id uint64
}

func (il *internalLogger) Info(v ...interface{}) {
	il.log(il.l.Info, v...)
}

func (il *internalLogger) Warning(v ...interface{}) {
	il.log(il.l.Warning, v...)
}

func (il *internalLogger) Error(v ...interface{}) {
	il.log(il.l.Error, v...)
}

func (il *internalLogger) Debug(v ...interface{}) {
	il.log(il.l.Debug, v...)
}

func (il *internalLogger) SetFileDepth(depth int) {
	if l, ok:=il.l.(setCallDepth); ok {
		if depth == 0 {
			depth = default_depth
		}
		l.SetFileDepth(depth)
	}
}

func (il *internalLogger) log(f func(v...interface{}), v ...interface{}) {
	if len(v) > 1 {
		v[0] = fmt.Sprintf("[%x] %s", il.id, v[0].(string))
	} else {
		v[0] = fmt.Sprintf("[%x] %v", il.id, v[0])
	}

	f(v...)
}


type setCallDepth interface {
	SetFileDepth(depth int)
}

type internalServerLoggerWriter struct {
	l Logger
}

func (islw *internalServerLoggerWriter) Write(b []byte) (n int, err error) {
	islw.l.Error("Go Server: %s", string(b))
	return len(b), nil
}