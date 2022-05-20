package rest

import (
	"fmt"
	"runtime"
)

var _handler = func(value any) {
	fmt.Println(value)
}

func SetHandler(handler func(any)) func(any) {
	_handler, handler = handler, _handler
	return handler
}

func WithRecover(togo func()) {
	defer Recover()
	togo()
}

func Recover() {
	if v := recover(); v != nil {
		_handler(v)
	}
}

func Panic(err error) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		panic(fmt.Errorf("%v(%v:%v): %w", runtime.FuncForPC(pc).Name(), file, line, err))
	} else {
		panic(err)
	}
}
