package chio

import "sync"

func Layer[I, O any](do func(I) O, workers int) func(<-chan I) <-chan O {
	return func(in <-chan I) <-chan O {
		var wg sync.WaitGroup
		out := make(chan O)
		wg.Add(workers)
		for ; workers > 0; workers-- {
			go func() {
				for v := range in {
					out <- do(v)
				}
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}
}

func FanIn[O any](done <-chan struct{}, fns ...func() O) (data <-chan O) {
	var wg sync.WaitGroup
	ch := make(chan O)
	do := func(fn func() O) {
		select {
		case ch <- fn():
		case <-done:
		}
		wg.Done()
	}
	wg.Add(len(fns))
	for _, fn := range fns {
		go do(fn)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
