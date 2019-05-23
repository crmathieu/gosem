# gosem
A semaphore implementation in golang with examples

Semaphores are tools that help manage concurrent accesses to a common resource (usually a buffer or a more complex data structure). In order to achieve this, a semaphore has a counter that indicates the level of availability the resource has at a given time.

A semaphore is a data structure using a channel as a mean to implement goroutine synchronization. The channel has a dimension (its capacity) that corresponds to the dimension of the sharable resource, and an initial count corresponding to its initial availability.

### Using gosem in your project
```go
import (
  sem "github.com/crmathieu/gosem/semaphore"
)
```

Imagine we want to share a buffer of 512 integers between 1 producer and 1 consumer. The producer wants to write to the buffer and the consumer wants to read from it.

```go
var buffer [512]int
var tail, head int = 0, 0
```

In order to protect this buffer and synchronize read and write operations, we are going to need 2 semaphores: one for reads and one for writes. The read semaphore is used to find out if there is anything to read from the buffer. The write semaphore is used to find out if there is any space available in the buffer so that we can write into it.

```go
readsem =  sem.Createsem("readsem",  512, 0)
writesem = sem.Createsem("writesem", 512, 512)
```

Note that both readsem and writesem have the same dimension: 512, but readsem has an initial value of 0 (because initially there is nothing to read) and writesem has an initial value of 512 (because initially the whole buffer is available).

you may also use the _CreateReadSemaphore_ or _CreateWriteSemaphore_ that will abstract the initial value given to the semaphore:

```go
readsem =  sem.CreateReadSemaphore("readsem",   512)
writesem = sem.CreateWriteSemaphore("writesem", 512)
```

The code of the producer looks like this:

```go
func producer() {
  i := 0
  for {
    writesem.Wait()
    buffer[head] = i
    i++
    head = (head+1) % 512
    readsem.Signal()
  }
}
```
In its producing loop, the producer first makes sure there is space in the buffer by calling writesem.Wait(). This call will return immediately if space is available but will block if the buffer is full. In the latter case, the call will return only after a consumer goroutine reads an entry from the buffer and performs a writesem.Signal() call to signify that one entry is now available.

Similarly, once a value was written in the buffer, the producer calls readsem.Signal() to indicate to the consumer that a value is available for consumption.

The code of the consumer looks like that:

```go
func consumer() {
  for {
    readsem.Wait()
    item = buffer[tail]
    tail = (tail+1) % 512
    writesem.Signal()
  }
}
```

In its consuming loop, the consumer first makes sure there is something to read from the buffer by calling readsem.Wait(). This call will return immediately if data is available but will block if the buffer is empty. In the latter case, the call will return only after a producer goroutine writes an entry to the buffer and performs a readsem.Signal() call to signify that one entry is now available for reading.

Similarly, once a value has been read from the buffer, the consumer calls writesem.Signal() to indicate to the producer that there is space available.


## Multiple consumers and producers
If we want to synchronize several consumers and producers accessing the same buffer, the code for both consumers and producers needs to handle concurrent access to the <b>head</b> and <b>tail</b> buffer indexes, because their value can be updated by multiple goroutines simultaneously (which was not the case in the previous example).

In order to do that, goroutines will need to have an exclusive access to these indexes when they update them. This is accomplished with the use of <b>mutex semaphores</b>. A mutex semaphore is like a normal semaphore with a capacity and an initial count of 1:

```go
mutex = sem.Createmutex("mymutex")
```
We are going to need a mutex to protect the <b>head</b> index and another mutex to protect the <b>tail</b> index:

```go
headmutex = sem.Createmutex("head-mutex")
tailmutex = sem.Createmutex("tail-mutex")
```

The producer code becomes:

```go
func producer() {
  i := 0
  for {
    writesem.Wait()
    headmutex.Enter()
    buffer[head] = i
    i++
    head = (head+1) % 512
    headmutex.Leave()
    readsem.Signal()
  }
}
```

and the consumer code becomes:

```go
func consumer() {
  for {
    readsem.Wait()
    tailmutex.Enter()
    item = buffer[tail]
    tail = (tail+1) % 512
    tailmutex.Leave()
    writesem.Signal()
  }
}
```
### Variable declaration
To declare a semaphore or a mutex:
var mysem *semaphore.Sem
-or-
var mymutex *semaphore.Mutex

### Semaphore API

First, import the gosem package:

```go
import (
  sem "github.com/crmathieu/gosem/semaphore"
)
```

#### Createsem: creates a counter semaphore

To create a semaphore with a capacity of 64, and an initial count of 0:
```go
sem1 = sem.Createsem("mySemaphore-1", 64, 0)
```
-or- to create a semaphore with a capacity of 64, and an initial count of 64:
```go
sem2 = sem.Createsem("mySemaphore-2", 64, 64)
```

#### Createmutex: creates a mutex  
```go
mutex = sem.Createmutex("myMutex")
```

Following a semaphore creation, there are a certain number of methods available to manipulate semaphores:

#### Reset
```go
mysem.Reset()

-or-

mymutex.Reset()
```

This will flush the semaphore internal channel and resets its counter to its original value.

#### Signal -or- V (-or- Leave)
```go
mysem.Signal()

-or-

mysem.V()
```
-or- for a mutex

```go
mymutex.Leave()
```
<b>Signal</b> and <b>V</b> accomplish the same thing which is to increase by 1 the level of availability of the resource. <b>Leave</b> is identical but reserved for <i>mutex</i>.

#### Wait -or- P (-or- Enter)
```go
mysem.Wait()

-or-

mysem.P()
```

-or- for a mutex

```go
mymutex.Enter()
```
<b>Wait</b> and <b>P</b> accomplish the same thing which is to decrease by 1 the level of availability of the resource. <b>Enter</b> is identical but reserved for <i>mutex</i>. When the semaphore counter reaches 0, the resource is no longer available, until a Signal (-or- a V) call is made by another goroutine.

#### Notes:
- The <b>P</b> / <b>V</b> notation comes from <b>Edsger Dijkstra</b>, who introduced the concept of semaphores in 1963. The letters are from the Dutch words <b>Probeer</b> (try) and <b>Verhoog</b> (increment).

- The terms <b>Enter</b> and <b>Leave</b> for a mutex refer to ```Entering``` and ```Leaving``` critical sections in your code. A <b>critical section</b> is a region in your code that can be executed only by one goroutine at a time. Typically, you will need to define a critical section everytime you need to access a resource that can potentially be modified by multiple goroutines. Once in the critical section, a goroutine is guaranteed to have exclusive access to the shared resource.

