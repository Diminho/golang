package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	defer timeTrack(time.Now(), "main")
	flag.Parse()
	inputJson := []byte(`[
        "https://d1ohg4ss876yi2.cloudfront.net/golang-resize-image/big.jpg",
        "http://www.textfiles.com/100/914bbs.txt",
        "http://www.textfiles.com/anarchy/001.txt",
        "http://www.textfiles.com/anarchy/ciaman.txt",
        "http://www.textfiles.com/anarchy/build_an_h_bomb.txt",
        "http://www.textfiles.com/anarchy/handbook.txt",
        "http://www.textfiles.com/anarchy/rocket.txt",
        "http://www.textfiles.com/computers/aboutems.txt",
        "http://www.textfiles.com/computers/act-13.txt",
        "http://www.textfiles.com/computers/ami-chts.txt",
        "http://www.textfiles.com/computers/amihist.txt",
        "http://www.textfiles.com/computers/arthayes.txt",
        "http://www.textfiles.com/computers/arcsuit.txt",
        "http://www.textfiles.com/fun/amtrak.txt",
        "http://www.textfiles.com/fun/amtrak1.txt",
        "http://www.textfiles.com/fun/celebrity.txt",
        "http://www.textfiles.com/fun/divright.txt",
        "http://www.textfiles.com/fun/faq11109.txt"

	]`)
	var urls []string

	if err := json.Unmarshal(inputJson, &urls); err != nil {
		panic(err)
	}

	limitation := flag.Arg(0)
	capacity, _ := strconv.Atoi(limitation)

	channel := make(chan string, capacity)

	for _, url := range urls {
		wg.Add(1)
		go download(url, channel)

	}
	wg.Wait()

}

func download(rawURL string, channel chan string) {
	channel <- rawURL
	fileURL, err := url.Parse(rawURL)
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	resp, err := http.Get(rawURL) // add a filter to check redirect

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer func() { <-channel }()

	fmt.Printf("%s with %v bytes downloaded", fileName, size)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
