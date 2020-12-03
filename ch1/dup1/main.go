package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func inefficient() {
	fmt.Println("inefficient")
	start := time.Now()
	s := os.Args[1]
	for _, args := range os.Args[2:] {
		s = s + " " + args
	}
	fmt.Println(s)
	mill := time.Since(start).Nanoseconds()
	fmt.Printf("elapsed time: %d\n", mill)
}

func efficient() {
	fmt.Println("efficient")
	start := time.Now()
	s := strings.Join(os.Args[1:], " ")
	fmt.Println(s)
	mill := time.Since(start).Nanoseconds()
	fmt.Printf("elapsed time: %d\n", mill)
}

func main() {
	efficient()
	inefficient()
	efficient()
	inefficient()
	efficient()
	inefficient()
}
