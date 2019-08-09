package main

/*
In broadcast, 1 or more producers feed a buffer and 1 or more consumers
read from the buffer. The difference between broadcast and xprodxcons is
that in broadcast, every consumer read the whole data, whereas in xprodxcons
each consumer reads a portion of the data.
*/
import (
	"fmt"
	"sync"
	"time"
	sem "github.com/crmathieu/gosem/pkg/semaphore"
	"github.com/crmathieu/gosem/pkg/queue"
)

var items, spaces *sem.Sem
var qmutex, pmutex, cmutex *sem.Mutex

const BUFFER_SIZE = 8192
var buffer [BUFFER_SIZE]int
var head, tail = 0, 0
var wg sync.WaitGroup
var qu *queue.QU

// main------------------------------------------------------------------------
func main() {

	pmutex = sem.Createmutex("producermutex")
	cmutex = sem.Createmutex("consumermutex")
	items = sem.Createsem("usedcount", BUFFER_SIZE, 0)
	spaces = sem.Createsem("availablecount", BUFFER_SIZE, BUFFER_SIZE)

	go broadcaster()

	// make sure to wait a bit to give a chance to
	// the broadcaster to initialize properly ...
	time.Sleep(10 * time.Millisecond)

	go producer(1000)
	go producer(1)

	for i:=0; i<13; i++ {
		go consumer(fmt.Sprintf("%c", 0x41+i))
	}

	// add more consumers to broadcast too after 10 secs
	time.Sleep(10 * time.Second)

	/*for j:=13; j<26; j++ {
		go consumer(fmt.Sprintf("%c", 0x41+j))
		time.Sleep(1 * time.Millisecond)
	}*/

	// wait a bit to give a chance to the goroutines to start
	time.Sleep(1 * time.Second)
	wg.Wait()
}

// producer--------------------------------------------------------------------
func producer(tag int) {
	defer wg.Done()
	//var tag int

	wg.Add(1)
	//for i:=0; i<BUFFER_SIZE; i++ {
	for {
		spaces.Wait()
		pmutex.Enter()
		produceItem(tag)
		head = (head + 1) % BUFFER_SIZE
		pmutex.Leave()
		items.Signal()
		tag++
	}
	fmt.Printf("Producer with tag = %d finished!\n", tag)
	fmt.Printf("Buffer = %v\n", buffer)
}

func produceItem(item int) {
	buffer[head] = item
}

// broadcaster-----------------------------------------------------------------
//            will fan out the data read from producer(s) to subscribers
// ----------------------------------------------------------------------------
func broadcaster()  {
	defer wg.Done()

	wg.Add(1)
	var item int
	qu = queue.InitQueue(BUFFER_SIZE)

	//for i:=0; i<BUFFER_SIZE; i++ {
	for {
		items.Wait()
		item = buffer[tail]
		qu.Fanout(item)
		tail = (tail + 1) % BUFFER_SIZE
		spaces.Signal()
	}
	fmt.Printf("broadcast finished!\n")
}

// consumer--------------------------------------------------------------------
func consumer(name string)  {
	defer wg.Done()

	wg.Add(1)
	var item, index int
	sq := qu.Subscribe()

	//for i:=0; i<BUFFER_SIZE; i++ {
	for {
		item, index = sq.Read()
		consumeItem(name, index, item)
	}

	fmt.Printf("Consumer finished!\n")
	qu.Close(sq)
}

func consumeItem(name string, index int, item int) {
	fmt.Printf("%s-%02d -> %d\n", name, index, item)
}
