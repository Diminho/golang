package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type RespFormat struct {
	Title  string `json:"title"`
	Number int    `json:"lines"`
}

func CountLines(file io.Reader) int {
	inputScanner := bufio.NewScanner(file)
	var lines int
	for inputScanner.Scan() {
		lines++
	}
	return lines
}

func FileToBeCount(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := CountLines(file)
	return lines
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		result := strings.Split(r.URL.Path, "/")
		bookName := result[len(result)-1]
		lines := FileToBeCount(os.Args[1] + bookName)
		data := RespFormat{bookName, lines}

		jsonData, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	})
	log.Fatal(http.ListenAndServe(":8081", nil))

}
