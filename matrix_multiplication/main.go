package main

import (
	"math/rand"
	"sync"
)

const (
	matrixSize = 200
)

var (
	matrixA   = [matrixSize][matrixSize]int{}
	matrixB   = [matrixSize][matrixSize]int{}
	result    = [matrixSize][matrixSize]int{}
	rwLock    = sync.RWMutex{}
	cond      = sync.NewCond(rwLock.RLocker())
	waitGroup = sync.WaitGroup{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5 // from -5 to 5
		}
	}
}

func calcRow(row int) {
	// 4: gets read lock
	rwLock.RLock()
	for {
		// 5: every of 200 goroutines subtracts 1 from the wait group
		waitGroup.Done()
		// 6: waits and releases the lock
		cond.Wait() // 13: continues to calculate the result
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}
}

func main() {
	// 1: adds 200 to the wait group
	waitGroup.Add(matrixSize)
	for row := 0; row < matrixSize; row++ {
		// 2: calls 200 goroutines with every row in the matrix
		go calcRow(row)
	}

	for i := 0; i < 100; i++ {
		// 3: waits
		waitGroup.Wait() // 7: continues
		// 8: gets write lock
		rwLock.Lock()
		// 9: generates 2 matrices
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		// 10: adds 200 to the wait group to wait on the next iteration
		waitGroup.Add(matrixSize)
		// 11: releases write lock
		rwLock.Unlock()
		// 12: notifies goroutines to continue
		cond.Broadcast()
	}
}
