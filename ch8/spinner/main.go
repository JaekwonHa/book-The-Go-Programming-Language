package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)

	const n = 45
	fibN := fibRecursive(n)
	//fibN := fibCache(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fibRecursive(x int) int {
	if x < 2 {
		return x
	}
	return fibRecursive(x-1) + fibRecursive(x-2)
}

var cache [46]int

func fibCache(x int) int {
	if x < 2 {
		return x
	}

	if cache[x] > 0 {
		return cache[x]
	}

	cache[x] = fibCache(x-1) + fibCache(x-2)
	return cache[x]
}
