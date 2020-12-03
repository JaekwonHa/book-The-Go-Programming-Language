package main

import (
	"fmt"
	"os"
)

func main() {
	for i, args := range os.Args {
		fmt.Printf("%d %s\n", i, args)
	}
}
