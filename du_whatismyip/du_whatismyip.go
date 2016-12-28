package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://httpbin.org/ip")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	result := make(map[string]string)
	dec := json.NewDecoder(resp.Body)
	errr := dec.Decode(&result)
	if errr != nil {
		log.Fatal(errr)
	}
	fmt.Println("My IP assress is:", result["origin"])

}
