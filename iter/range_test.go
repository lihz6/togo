package iter_test

import (
	"fmt"
	"strings"
)

func Example_string_rune() {
	str := "Hi世界"
	for i, s := range str {
		fmt.Println(i, string(s))
	}
	// Output:
	// 0 H
	// 1 i
	// 2 世
	// 5 界
}

func Example_string_byte() {
	str, buf := "Hi世界", strings.Builder{}
	// for i := 0; i < len(str); i++ {
	for i, b := range []byte(str) {
		if i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(fmt.Sprintf("%x", b))
	}
	fmt.Println(buf.String())
	// Output: 48 69 e4 b8 96 e7 95 8c
}

func Example_array() {
	arr := [...]int{1, 2, 3}
	for i, a := range arr[:] {
		arr[i] = a * a
	}
	fmt.Println(arr[:])
	// Output: [1 4 9]
}

func Example_slice() {
	arr := []int{1, 2, 3}
	for i, a := range arr {
		arr[i] = a * a
	}
	fmt.Println(arr)
	// Output: [1 4 9]
}

func Example_map() {
	m := map[int]int{0: 1, 1: 2, 2: 3}
	for k, v := range m {
		fmt.Println(k, v)
	}
	// Unordered Output:
	// 0 1
	// 1 2
	// 2 3
}

func Example_channel() {
	ch := make(chan int)
	go func() {
		defer close(ch)
		ch <- 1
		ch <- 2
		ch <- 3
	}()
	for n := range ch {
		fmt.Println(n)
	}
	// Output:
	// 1
	// 2
	// 3
}
