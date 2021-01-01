package main

import "fmt"

func main() {
	fmt.Println(mirroredQuery())
}

func mirroredQuery() string {
	responses := make(chan string, 3)
	go func() { responses <- request("asia.gopl.io") }()
	go func() { responses <- request("europe.gopl.io") }()
	go func() { responses <- request("americas.gopl.io") }()
	return <-responses // 가장 빠른 응답 반환
}

func request(response string) string {
	return response
}
