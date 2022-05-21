package iter_test

import (
	"errors"
	"fmt"
	"testing"
	"togo/iter"
)

func Example_channel_break() {
	done := make(chan struct{})
	next := iter.Chan(done, fibClosure())
	defer close(done)
	for n := range next {
		fmt.Print(n)
		if n > 10 {
			break
		}
	}
	// Output: 11235813
}

func Example_channel_send() {
	done := make(chan struct{})
	next := iter.Chan(done, fibClosure())
	defer close(done)
	for n := range next {
		fmt.Print(n)
		if n > 10 {
			done <- struct{}{}
		}
	}
	// Output: 11235813
}

func TestFibChannel_close(t *testing.T) {
	done := make(chan struct{})
	next := iter.Chan(done, fibClosure())
	close(done)
	<-next
	if v, ok := <-next; ok {
		t.Errorf("expect next closed, but not: %v", v)
	}
}

func TestFibChan_send(t *testing.T) {
	done := make(chan struct{})
	next := iter.Chan(done, fibClosure())
	done <- struct{}{}
	if v, ok := <-next; ok {
		t.Errorf("expect next closed, but not: %v", v)
	}
	// it's OK to leave the done open
}

func Example_error() {
	iterable := func(n int) func(next *int, err *error) (ok bool) {
		i := -1
		return func(next *int, err *error) (ok bool) {
			i++
			if i > n {
				*err = errors.New("iterator exceeded")
				return false
			}
			*next = i
			return i < n
		}
	}
	var err error
	iter := iterable(N)
	for next := 0; iter(&next, &err); {
	}
	iter(nil, &err)
	if err != nil {
		fmt.Println(err)
	}
	// Output: iterator exceeded
}

func Example_closure() {
	iter := fibClosure()
	for n := 0; iter(&n); {
		fmt.Print(n)
		if n > 10 {
			break
		}
	}
	// Output: 11235813
}

func fibClosure() iter.Iter[int] {
	a, b := 0, 1
	return func(next *int) bool {
		a, b = b, a+b
		*next = a
		return true
	}
}
