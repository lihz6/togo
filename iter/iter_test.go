package iter_test

import "testing"

const N = 100_000_000

func BenchmarkNextOk(b *testing.B) {
	iterable := func(n int) func(next *int) (ok bool) {
		i := -1
		return func(next *int) (ok bool) {
			i++
			*next = i
			return i < n
		}
	}
	for i := 0; i < b.N; i++ {
		iter := iterable(N)
		sum := 0
		for next := 0; iter(&next); {
			sum += next
		}
	}
}

func BenchmarkNextErrOk(b *testing.B) {
	iterable := func(n int) func(next *int, err *error) (ok bool) {
		i := -1
		return func(next *int, err *error) (ok bool) {
			i++
			*next = i
			return i < n
		}
	}
	var err error
	for i := 0; i < b.N; i++ {
		iter := iterable(N)
		sum := 0
		for next := 0; iter(&next, &err); {
			sum += next
		}
	}
}

func BenchmarkNextOkNew(b *testing.B) {
	iterable := func(n int) func(next *int) (ok bool) {
		i := -1
		return func(next *int) (ok bool) {
			i++
			*next = i
			return i < n
		}
	}
	for i := 0; i < b.N; i++ {
		iter := iterable(N)
		sum := 0
		for next := new(int); iter(next); {
			sum += *next
		}
	}
}

func BenchmarkNextDone(b *testing.B) {
	iterable := func(n int) func() (next int, done bool) {
		i := -1
		return func() (next int, done bool) {
			i++
			return i, i >= n
		}
	}
	for i := 0; i < b.N; i++ {
		iter := iterable(N)
		sum := 0
		for {
			next, done := iter()
			if done {
				break
			}
			sum += next
		}
	}
}

func BenchmarkNextSkip(b *testing.B) {
	iterable := func(n int) func(func(next int) (skip bool)) {
		return func(f func(next int) (skip bool)) {
			for i := 0; i < n; i++ {
				if f(i) {
					return
				}
			}
		}
	}
	for i := 0; i < b.N; i++ {
		iter := iterable(N)
		sum := 0
		iter(func(next int) (skip bool) {
			sum += next
			return
		})
	}
}

func BenchmarkHasNextNext(b *testing.B) {
	iterable := func(n int) (hasNext func() bool, next func() int) {
		i := 0
		hasNext = func() bool {
			return i < n
		}
		next = func() int {
			n := i
			i++
			return n
		}
		return
	}
	for i := 0; i < b.N; i++ {
		hasNext, next := iterable(N)
		sum := 0
		for hasNext() {
			sum += next()
		}

	}
}
