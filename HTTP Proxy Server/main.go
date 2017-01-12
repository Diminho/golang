package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			// w.Write([]byte("IT WAS GET METHOD"))

			resp, err := http.Get(r.RequestURI)
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			w.Write(body)
		}
	default:
		{
			log.Print("Cannot handle method ", r.Method)
			http.Error(w, "Only GET method", http.StatusNotImplemented)
			return
		}

	}
}

func main() {

	http.HandleFunc("/", proxyHandler)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
