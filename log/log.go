package log

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(any)
}

type PrintLogger struct {
	queue chan any
	wait  chan any
}

func NewPrintLogger() *PrintLogger {
	return &PrintLogger{make(chan any, 256), make(chan any, 1)}
}

func (l PrintLogger) Log(log any) {
	l.queue <- log
}

func (l PrintLogger) Run() {
	for log := range l.queue {
		fmt.Println(timestamp() + " " + fmt.Sprint(log))
	}
	l.wait <- nil
}

func (l PrintLogger) Wait() {
	<-l.wait
}

func (l PrintLogger) Shutdown() {
	close(l.queue)
}

func timestamp() string {
	return time.Now().UTC().Format("[2006-01-02T15:04:05UTC]")
}
