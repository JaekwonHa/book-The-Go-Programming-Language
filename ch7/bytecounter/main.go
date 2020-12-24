package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// 7.1
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

// ex7.1
type WordCounter int

func (c *WordCounter) Write(p []byte) (count int, e error) {
	count = retCount(p, bufio.ScanWords)
	*c += WordCounter(count)
	return count, nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (count int, e error) {
	count = retCount(p, bufio.ScanLines)
	*c += LineCounter(count)
	return count, nil
}

func retCount(p []byte, fn bufio.SplitFunc) (count int) {
	s := string(p)
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(fn)
	count = 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	return
}

// ex7.2
type WrapperWriter struct {
	w       io.Writer
	written int64
}

func (c *WrapperWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.written += int64(n)
	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &WrapperWriter{w, 0}
	return c, &c.written
}

func main() {
	var c ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)

	var c2 WordCounter
	c2.Write([]byte("hello"))
	fmt.Fprintf(&c2, "hello, %s", name)
	fmt.Println(c2)

	var c3 LineCounter
	c3.Write([]byte("hello"))
	fmt.Fprintf(&c3, "hello, %s\nhello", name)
	fmt.Println(c3)

	c4, n := CountingWriter(&c)
	c4.Write([]byte("hello"))
	fmt.Printf("%d\n", *n)

	root := &tree{value: 3}
	root = add(root, 4)
	root = add(root, 1)
	fmt.Println(root)
}
