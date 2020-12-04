package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	var result string
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		result += <-ch
	}
	result += fmt.Sprintf("%.2fs elapsed\n", time.Since(start).Seconds())
	fmt.Printf(result)

	uuid, _ := exec.Command("uuidgen").Output()
	ioutil.WriteFile("./fetchall_result_"+string(uuid), []byte(result), 0644)
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s\n", secs, nbytes, url)
}
