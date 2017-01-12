package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

type conn interface {
	doQuery()
}

type connection struct {
	httpWrap http.Client
}

func main() {

	resources := []string{
		"http://football.ua/",
		"http://www.juventus.com/en/",
		"http://en.valenciacf.com/",
		"http://www.liverpoolecho.co.uk/",
		"http://www.bvb.de/eng/",
		"http://www.psg.fr/en/Accueil/0/Home"}

	var alternative http.Client
	structure := connection{alternative}
	// structure.doQuery("http://www.liverpoolecho.co.uk/")
	fmt.Printf("%s", query(structure, resources))

}

func query(conn connection, queries []string) string {
	ch := make(chan string, 1)
	for _, query := range queries {
		go func(query string) {
			select {
			case ch <- conn.doQuery(query):
			default:
			}
		}(query)
	}
	return <-ch
}

func (c connection) doQuery(query string) string {
	resp, err := c.httpWrap.Get(query)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)

	finalString := string(body)
	fmt.Printf("%s\n", query)
	return finalString

}
