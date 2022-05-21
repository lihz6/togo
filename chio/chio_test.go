package chio_test

import (
	"math/rand"
	"testing"
	"time"
	"togo/chio"
)

func TestLayer(t *testing.T) {
	scale := rand.Intn(1000)
	done, next := produce(scale)
	defer close(done)
	sq := chio.Layer(sleep(square, time.Second), scale)
	in := chio.Layer(inverse, 1)
	ni := chio.Layer(inverse, 1)
	sum := 0
	for n := range ni(in(sq(next))) {
		sum += n
	}
	if sum != sumUp(scale) {
		t.Errorf("want %v, got %v", sumUp(scale), sum)
	}
}

func BenchmarkLayer_1000_1(b *testing.B) {
	const scale = 1000
	done, next := produce(scale)
	defer close(done)
	sq := chio.Layer(sleep(square, time.Second), scale)
	in := chio.Layer(inverse, 1)
	ni := chio.Layer(inverse, 1)
	for i := 0; i < b.N; i++ {
		sum := 0
		for n := range ni(in(sq(next))) {
			sum += n
		}
		if sum != sumUp(scale) {
			b.Errorf("want %v, got %v", sumUp(scale), sum)
		}
	}
}

func BenchmarkLayer_1_1(b *testing.B) {
	const scale = 8
	done, next := produce(scale)
	defer close(done)
	sq := chio.Layer(sleep(square, time.Second), 1)
	in := chio.Layer(inverse, 1)
	ni := chio.Layer(inverse, 1)
	for i := 0; i < b.N; i++ {
		sum := 0
		for n := range ni(in(sq(next))) {
			sum += n
		}
		if sum != sumUp(scale) {
			b.Errorf("want %v, got %v", sumUp(scale), sum)
		}
	}
}

func square(i int) int {
	return i * i
}

func inverse(i int) int {
	return -i
}

func produce(n int) (chan<- struct{}, <-chan int) {
	done, next := make(chan struct{}), make(chan int)
	go func() {
		defer close(next)
		for i := 0; i < n; i++ {
			select {
			case next <- i:
			case <-done:
				return
			}
		}
	}()
	return done, next
}

func sleep[I, O any](fn func(I) O, duration time.Duration) func(I) O {
	return func(i I) O {
		time.Sleep(duration)
		return fn(i)
	}
}

func sumUp(scale int) (sum int) {
	done, next := produce(scale)
	defer close(done)
	for n := range next {
		sum += inverse(inverse(square(n)))
	}
	return
}
