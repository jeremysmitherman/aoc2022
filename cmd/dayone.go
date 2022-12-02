package main

import (
	"aoc2022/dayone"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Ensure we get a file name as the argument
	if len(os.Args) != 2 {
		fmt.Printf("Usage: ./dayone [INPUT FILE]")
		os.Exit(1)
	}

	// Open the input file, close the handle on exit
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("File not found: %s\n", os.Args[1])
	}
	defer f.Close()

	// Create our sorter and counter builder closure.  Create a waitgroup for the goroutines
	calorieSorter := dayone.NewCounterSorter()
	buildCalorieCounter := dayone.NewCalorieCounter()
	wg := &sync.WaitGroup{}

	// Build a scanner to read our input file easily and set it to split on newlines
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// Build the initial calorie counter and have it start listening in a goroutine
	calorieCounter := buildCalorieCounter()
	go calorieCounter.Listen(wg)

	for {
		for scanner.Scan() {
			input := scanner.Text()
			// Blank line, so close the input channel on the current counter, start a goroutine to insert it into the
			// calorie sorter asynchronously, then create a new calorie counter and continue scanning
			if input == "" {
				close(calorieCounter.Input)
				go calorieSorter.Insert(calorieCounter, wg)

				calorieCounter = buildCalorieCounter()
				go calorieCounter.Listen(wg)
				continue
			} else {
				// Value found on this line, convert the input, send it to the input channel of the current calorie counter
				val, err := strconv.Atoi(input)
				if err != nil {
					log.Fatalf(err.Error())
				}
				calorieCounter.Input <- val
			}
		}

		// End of the file.  Close the current counter, and asynchronously send it to the sorter then end the loop.
		close(calorieCounter.Input)
		go calorieSorter.Insert(calorieCounter, wg)
		break
	}

	// Wait for all the counters and the sorter to finish, Get the last entry of the sorter's ordered list and print the
	// calorie total it contains.
	wg.Wait()
	fmt.Println(calorieSorter.OrderedList[len(calorieSorter.OrderedList)-1].Total)

	top3Sum := 0
	for i := 1; i <= 3; i++ {
		v := calorieSorter.OrderedList[len(calorieSorter.OrderedList)-i].Total
		top3Sum += v
	}

	fmt.Println(top3Sum)

}
