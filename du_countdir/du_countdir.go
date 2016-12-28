package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

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

func ProcessDir(dir string) {
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if ext := filepath.Ext(file.Name()); ext == ".txt" {
			lines := FileToBeCount(dir + file.Name())
			fmt.Printf("%s\t%d\n", file.Name(), lines)
		}
	}
}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		ProcessDir(arg)
	}

}
