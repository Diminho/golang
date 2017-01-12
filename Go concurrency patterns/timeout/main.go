package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seed := rand.NewSource(22) // new seed
		newRandom := rand.New(seed)
		time.Sleep(time.Second * time.Duration(newRandom.Intn(5)))
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	fmt.Println(ts.URL)

	response := performRequest(ts.URL)
	fmt.Println(string(response))

}

func performRequest(url string) string {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err.Error()
	}
	finalString := string(greeting)
	return finalString

}
