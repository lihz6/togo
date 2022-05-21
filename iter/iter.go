package iter

import "time"

type IterErr[T any] func(next *T, err *error) (ok bool)

type Iter[T any] func(next *T) (ok bool)

func Chan[T any](done <-chan struct{}, iter func(*T) bool) (next <-chan T) {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for n := new(T); iter(n); {
			select {
			case ch <- *n:
			case <-done:
				return
			}
		}
	}()
	return ch
}

func Timeout[T any](d time.Duration, ch <-chan T) Iter[T] {
	timer := time.NewTimer(d)
	return func(next *T) bool {
		select {
		case <-timer.C:
			return false
		case v, ok := <-ch:
			if !timer.Stop() {
				<-timer.C
			}
			if ok {
				timer.Reset(d)
			}
			*next = v
			return ok
		}
	}
}
