package logs_test

import (
	"time"
	"togo/logs"
)

func Example() {
	defer logs.Duration(time.Now())
	time.Sleep(time.Second)
}
