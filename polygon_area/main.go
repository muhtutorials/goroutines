package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point2D struct {
	x, y int
}

const numberOfThreads int = 4

var r = regexp.MustCompile(`\((\d*),(\d*)\)`)
var wg sync.WaitGroup

func findArea(inputChannel chan string) {
	defer wg.Done()
	for pointStr := range inputChannel {
		var points []Point2D
		for _, p := range r.FindAllStringSubmatch(pointStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x: x, y: y})
		}
		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)] // %len(points) returns i to zero
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
}

func main() {
	path := "./polygons.txt"
	data, _ := os.ReadFile(path)
	text := string(data)

	inputChannel := make(chan string, 10)
	wg.Add(numberOfThreads)
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChannel)
	}
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		inputChannel <- line
	}
	close(inputChannel)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}
