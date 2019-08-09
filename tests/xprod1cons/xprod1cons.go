package main

// testing synchronization of multiple producers and 1 consumer on a shared resource
import (
	"fmt"
	"sync"
	"time"
	sem "github.com/crmathieu/gosem/pkg/semaphore"
)

var items, spaces *sem.Sem
var pmutex *sem.Mutex

const BUFFER_SIZE = 32
var buffer [BUFFER_SIZE]int
var head, tail = 0, 0
var wg sync.WaitGroup

// main------------------------------------------------------------------------
func main() {

	pmutex = sem.Createmutex("producermutex")
	items = sem.Createsem("usedcount", BUFFER_SIZE, 0)
	spaces = sem.Createsem("availablecount", BUFFER_SIZE, BUFFER_SIZE)

	go producer(11111)
	go producer(11112)
	go producer(11113)
	go producer(11114)
	go producer(11115)
	go producer(11116)
	go producer(11117)
	go producer(11118)
	go producer(11119)
	go producer(11110)
	go producer(21111)
	go producer(21112)
	go producer(21113)
	go producer(21114)
	go producer(21115)
	go producer(21116)
	go producer(21117)
	go producer(21118)
	go producer(21119)
	go producer(21110)

	go consumer("A")

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
		item = buffer[tail]
		consumeItem(name, tail, item)
		tail = (tail + 1) % BUFFER_SIZE
		spaces.Signal()
	}
	fmt.Printf("Consumer finished!\n")
}

func consumeItem(name string, index int, item int) {
	fmt.Printf("%s%02d -> %d\n", name, index, item)
}
