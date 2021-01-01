package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)

	// 커맨드라인 인자로 시작
	go func() { worklist <- os.Args[1:] }()
	//worklist <- os.Args[1:]

	// 웹을 동시에 시작
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
