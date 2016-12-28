package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func CountLines(file io.Reader) int {
	inputScanner := bufio.NewScanner(file)
	var lines int
	for inputScanner.Scan() {
		lines++
	}
	return lines
}

func FileToBeCount(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := CountLines(file)
	fmt.Printf("%s\t%d\n", filepath, lines)

}

func main() {
	args := os.Args[1:]
	fmt.Println(os.Args)
	for _, arg := range args {
		FileToBeCount(arg)
	}
}
