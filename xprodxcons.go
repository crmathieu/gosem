package main

import (
	"fmt"
	"sync"
	"time"
	sem "github.com/crmathieu/gosem/semaphore"
)

var items, spaces *sem.Sem
var pmutex, cmutex *sem.Mutex

const BUFFER_SIZE = 32
var buffer [BUFFER_SIZE]int
var head, tail = 0, 0
var wg sync.WaitGroup

// main------------------------------------------------------------------------
func main() {

	pmutex = sem.Createmutex("producermutex")
	cmutex = sem.Createmutex("consumermutex")
	items = sem.Createsem("usedcount", BUFFER_SIZE, 0)
	spaces = sem.Createsem("availablecount", BUFFER_SIZE, BUFFER_SIZE)

	go producer(11111)
	go producer(22222)
	go producer(33333)
	go producer(44444)
	go producer(55555)
	go consumer("A")
	go consumer("B")
	go consumer("C")
	go consumer("D")

	// wait a bit to give a chance to the goroutines to start
	time.Sleep(1 * time.Second)
	wg.Wait()

}

// producer--------------------------------------------------------------------
func producer(tag int) {
	defer wg.Done()

	wg.Add(1)

	//for i := 0; i< BUFFER_SIZE; i++  {
	for {
		spaces.Wait()
		pmutex.Enter()
		produceItem(tag)
		head = (head + 1) % BUFFER_SIZE
		pmutex.Leave()
		items.Signal()
	}
	fmt.Printf("Producer with tag = %d finished!\n", tag)
}

func produceItem(item int) {
	buffer[head] = item
}

// consumer--------------------------------------------------------------------
func consumer(name string)  {
	defer wg.Done()

	wg.Add(1)
	var item int
  //for i := 0; i< BUFFER_SIZE*2; i++  {
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
