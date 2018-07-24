package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	routinesLimit := 5
	limiter := make(chan struct{}, routinesLimit)

	var total uint32
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		wg.Add(1)
		go counter(scanner.Text(), &total, limiter, &wg)
	}

	wg.Wait()
	fmt.Printf("Total: %d \n", total)
}

func counter(url string, total *uint32, limiter chan struct{}, wg *sync.WaitGroup) {
	limiter <- struct{}{}

	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	count := uint32(strings.Count(string(body), "Go"))
	fmt.Printf("Count for %s %d \n", url, count)
	atomic.AddUint32(total, count)

	<-limiter
	wg.Done()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
