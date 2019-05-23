package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/crmathieu/gosem/semaphore"
)

var items, spaces *semaphore.Sem
var cmutex *semaphore.Mutex

const BUFFER_SIZE = 32
var buffer [BUFFER_SIZE]int
var head, tail = 0, 0
var wg sync.WaitGroup

// main------------------------------------------------------------------------
func main() {

	cmutex = semaphore.Cmutex("consumermutex")
	items = semaphore.Csem("usedcount", BUFFER_SIZE, 0)
	spaces = semaphore.Csem("availablecount", BUFFER_SIZE, BUFFER_SIZE)

	go producer()
	go consumer("A")
	go consumer("B")
	go consumer("C")
	go consumer("D")
	go consumer("E")
	go consumer("F")
	go consumer("G")
	go consumer("H")
	go consumer("I")

	// wait a bit to give a chance to the goroutines to start
	time.Sleep(1 * time.Second)
	wg.Wait()

}

// producer--------------------------------------------------------------------
func producer() {
	defer wg.Done()
	var tag int
	wg.Add(1)
	for {
		spaces.Wait()
		produceItem(tag)
		head = (head + 1) % BUFFER_SIZE
		items.Signal()
		tag++
	}
	fmt.Printf("Producer with tag = %d finished!\n", tag)
	fmt.Printf("Buffer = %v\n", buffer)
}

func produceItem(item int) {
	buffer[head] = item
}

// consumer--------------------------------------------------------------------
func consumer(name string)  {
	defer wg.Done()

	wg.Add(1)
	var item int
	for {
		items.Wait()
		cmutex.Enter()
		item = buffer[tail]
		consumeItem(name, tail, item)
		tail = (tail + 1) % BUFFER_SIZE
		cmutex.Leave()
		spaces.Signal()
	}
	fmt.Printf("Consumer finished!\n")
}

func consumeItem(name string, index int, item int) {
	fmt.Printf("%s%02d -> %d\n", name, index, item)
}
