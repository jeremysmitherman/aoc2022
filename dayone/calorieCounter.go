package dayone

import "sync"

// CalorieCounter contains a channel to take input, and ID to identify the struct, and the total calorie count
type CalorieCounter struct {
	Input chan int
	ID    int
	Total int
}

// Listen takes an int input and adds it to the total until the Input channel is closed
func (c *CalorieCounter) Listen(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		v, open := <-c.Input
		if !open {
			wg.Done()
			break
		} else {
			c.Total += v
		}
	}
}

// NewCalorieCounter returns a function that will give you a new Calorie counter with a unique ID
func NewCalorieCounter() func() *CalorieCounter {
	nextID := 0
	return func() *CalorieCounter {
		nextID += 1
		return &CalorieCounter{
			Input: make(chan int),
			ID:    nextID,
		}
	}
}

// CounterSorter contains an ordered list of CalorieCounter sorted by the total calories
type CounterSorter struct {
	OrderedList []*CalorieCounter
	lock        *sync.Mutex
}

func NewCounterSorter() *CounterSorter {
	return &CounterSorter{
		lock: &sync.Mutex{},
	}
}

// Insert takes a counter and inserts it in order of Total calories
func (c *CounterSorter) Insert(counter *CalorieCounter, waitGroup *sync.WaitGroup) {
	// Add to the waitgroup and lock the mutex
	waitGroup.Add(1)
	c.lock.Lock()

	// Clear locks and inform WG that we're finished.
	defer func() {
		waitGroup.Done()
		c.lock.Unlock()
	}()

	// If we're the biggest known total, just add it to the end and call it a day
	if len(c.OrderedList) == 0 || counter.Total > c.OrderedList[len(c.OrderedList)-1].Total {
		c.OrderedList = append(c.OrderedList, counter)
		return
	}

	// I love Go sometimes
	var newList []*CalorieCounter
	for i, currentCounter := range c.OrderedList {
		if currentCounter.ID != counter.ID && counter.Total <= currentCounter.Total {
			newList = append(newList, counter)
			newList = append(newList, c.OrderedList[i:len(c.OrderedList)]...)
			break
		} else {
			newList = append(newList, currentCounter)
		}
	}
	c.OrderedList = newList
}
