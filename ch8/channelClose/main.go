package main

import "fmt"

func main() {

	c := make(chan int, 3)

	c <- 0

	c <- 0

	c <- 0

	close(c)

	for i := 0; i < 6; i++ {

		v, ok := <-c

		fmt.Printf("%d, %v \n", v, ok)

	}

}
