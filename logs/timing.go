package logs

import (
	"log"
	"runtime"
	"time"
)

func Duration(invocation time.Time) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		name := runtime.FuncForPC(pc).Name()
		log.Printf("[FunctionDuration] %v(%v:%v): %v", name, file, line, time.Since(invocation))
		return
	}
	log.Printf("[FunctionDuration] unknown(...): %v", time.Since(invocation))
}
