package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"net/http"
	"net/http/httptest"
	"time"
)

var saltInt int
var test string

func main() {
	saltInt = 100
	channel := make(chan string, 5)

	for i := 0; i < 5; i++ {
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		fmt.Println(ts.URL)
		go performRequest(ts.URL, channel)
	}

	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("timed out")
			return
		case msg1 := <-channel:
			fmt.Println("received", msg1)

		}
	}
	// response := performRequest(ts.URL, channel)

}

func handler(w http.ResponseWriter, r *http.Request) {

	segments := strings.Split(r.Host, ":")
	port := segments[len(segments)-1]
	int, _ := strconv.ParseInt(port, 10, 64)
	newRandom := rand.New(rand.NewSource(time.Now().UnixNano() + int)) // new seed
	randonInt := newRandom.Intn(5)
	fmt.Println("Random:", randonInt, "host: ", r.Host)
	time.Sleep(time.Second * time.Duration(randonInt))
	fmt.Fprintln(w, "Hello, client from ", r.Host)

}

func performRequest(url string, channel chan string) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	finalString := string(greeting)
	test = finalString

	// fmt.Println(finalString + url)
	channel <- url

}
