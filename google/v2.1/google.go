package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("Web")
	Image = fakeSearch("Image")
	Video = fakeSearch("Video")
)

type Search func(query string) string

func fakeSearch(kind string) Search {
	return func(query string) string {
		duration := rand.Intn(100)
		time.Sleep(time.Duration(duration) * time.Millisecond)
		return fmt.Sprintf("%s result for %q %d ms\n", kind, query, duration)
	}
}

func main() {
	flag.Parse()
	query := flag.Arg(0)

	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google(query)
	elapsed := time.Since(start)

	fmt.Println(results)
	fmt.Println(elapsed)
}

func Google(query string) (results []string) {
	c := make(chan string)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	// 3回結果を受け取る or 10ms経つ まで待つ
	timeout := time.After(10 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}
