package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyz")

const time_period int = 300
const no_workers = 16

var shakespeare = []byte("shakespeare")

// Compares the contents of the generated string consecutively and returns the highest match
// Unsafe, but we know the array is of a fixed size so this is fine. Remove bound checking == quicker comparisons

func Compare(generated []byte) int {
	shakespeare_compare := *(*[11]byte)(unsafe.Pointer(&shakespeare[0]))
	generated_compare := *(*[11]byte)(unsafe.Pointer(&generated[0]))

	for i := 0; i < 11; i++ {
		if generated_compare[i] != shakespeare_compare[i] {
			return i
		}
	}
	return 11
}

// Generates a random 11-letter string of [a-z]
func Generate() []byte {
	b := make([]byte, 11)
	for i := range b {
		b[i] = letters[rand.Intn(26)]
	}
	return b
}

// Worker function that continuously generates strings and updates max_matches if a better score is found
func Worker(done chan interface{}, results chan<- int, guesses chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	local_max := 0
	no_matches := 0
	keep_working := true
	for keep_working {
		select {
		case <-done:
			keep_working = false
		default:
			{
				no_matches++
				score := Compare(Generate())
				if score > local_max {
					local_max = score
				}
			}
		}
	}
	results <- local_max
	guesses <- no_matches
}

func main() {
	// Set up threads to concurrently execute
	fmt.Printf("Starting shakespeare with %d worker threads for %d seconds\n", no_workers, time_period)
	guesses := make(chan int, no_workers)
	matches := make(chan int, no_workers)
	done := make(chan interface{})
	wg := new(sync.WaitGroup)
	for i := 0; i < no_workers; i++ {
		wg.Add(1)
		go Worker(done, matches, guesses, wg)
	}
	// Run for time period
	time.Sleep(time.Second * time.Duration(time_period))
	close(done)
	wg.Wait()

	// Process channels to retrieve total guesses + best max
	no_guesses := 0
	max_matches := 0
	for i := 0; i < no_workers; i++ {
		max_matches = max(<-matches, max_matches)
		no_guesses = no_guesses + (<-guesses)
	}
	fmt.Printf("Number of guesses: %d\n", no_guesses)
	fmt.Printf("Best score: %d\n", max_matches)
}
