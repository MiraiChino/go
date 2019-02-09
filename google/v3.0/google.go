package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	Webs   = replicas("Web", 100)
	Images = replicas("Image", 100)
	Videos = replicas("Image", 100)
)

func replicas(kind string, n int) (replicas []Search) {
	for i := 1; i < n+1; i++ {
		server := fmt.Sprintf("%s %d", kind, i)
		replicas = append(replicas, fakeSearch(server))
	}
	return
}

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
	go func() { c <- First(query, Webs...) }()
	go func() { c <- First(query, Images...) }()
	go func() { c <- First(query, Videos...) }()

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

func First(query string, replicas ...Search) string {
	c := make(chan string)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}
