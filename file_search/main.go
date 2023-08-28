package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitGroup sync.WaitGroup
	lock      sync.Mutex
)

func searchFile(root, filename string) {
	fmt.Println("Searching in", root)
	files, _ := os.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			waitGroup.Add(1)
			go searchFile(filepath.Join(root, file.Name()), filename)
		}
	}
	waitGroup.Done()
}

func main() {
	waitGroup.Add(1)
	go searchFile("D:\\Downloads", "Redis-For-Dummies_FINAL.pdf")
	waitGroup.Wait()
	for _, file := range matches {
		fmt.Println("Matched", file)
	}
}
