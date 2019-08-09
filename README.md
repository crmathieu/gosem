# gosem

Semaphores are tools that help manage concurrent accesses to a common resource (the resource is usually one or more data structures with a set of indexes or pointers used to manipulate the data). In order to achieve this, a semaphore has a counter that indicates the level of availability the resource has at a given time.

In the context of Golang, a semaphore will be represented as a data structure using a channel as a mean to provide goroutines synchronization. The channel has a dimension (its capacity) that corresponds to the dimension of the sharable resource, and an initial count corresponding to its initial availability.

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
    i = (i + 1) % 4096
    head = (head+1) % 512
    readsem.Signal()
  }
}
```
In its loop, the **producer** first makes sure there is available space in the buffer by calling _writesem.Wait()_. This call will return immediately if space is available but will block if the buffer is full. In the latter case, the call will return only after the **consumer** goroutine reads an entry from the buffer and performs a _writesem.Signal()_ call to indicate that one entry is now available.

Similarly, once a value was written in the buffer, the **producer** calls _readsem.Signal()_ to indicate that one entry is available for consumption.

The code of the consumer looks like that:

```go
func consumer() {
  for {
    readsem.Wait()
    item = buffer[tail]
    tail = (tail+1) % 512
    writesem.Signal()
    fmt.Println(item)
  }
}
```

In its loop, the **consumer** first makes sure there is something to read from the buffer by calling _readsem.Wait()_. This call will return immediately if data is available but will block if the buffer is empty. In the latter case, the call will return only after the **producer** goroutine writes an entry to the buffer and performs a _readsem.Signal()_ call to indicate that one entry is ready to be consumed.

Similarly, once a value has been read from the buffer, the **consumer** calls _writesem.Signal()_ to indicate that space is available for production.


## Multiple consumers and producers
If we want to synchronize several consumers and producers accessing the same buffer, the code for both consumers and producers needs to handle concurrent access to the <b>head</b> and <b>tail</b> buffer indexes, because their value can be updated by multiple goroutines simultaneously (which was not the case in the previous example).

In order to do that, goroutines will need to have an exclusive access to these indexes when they update them. This is accomplished with the use of <b>mutex semaphores</b>. A mutex semaphore is like a normal semaphore with a capacity and an initial count of 1:

```go
mutex = sem.Createmutex("mymutex")
```
We are going to need a mutex to protect the <b>head</b> index used by multiple producers and another mutex to protect the <b>tail</b> index used by multiple consumers:

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
    i = (i + 1) % 4096    
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
    fmt.Println(item)
  }
}
```

### Semaphore API

First, import the gosem package:

```go
import (
  sem "github.com/crmathieu/gosem/pkg/semaphore"
)
```

#### Variable declaration
To declare a semaphore or a mutex:
```go
var mysem *semaphore.Sem
```
-or-
```go
var mymutex *semaphore.Mutex
```

#### Createsem: creates a counter semaphore
_func Createsem(name string, capacity int, initialcount int) *Sem_

To create a semaphore with a capacity of 64, and an initial count of 0:
```go
mysem := sem.Createsem("mySemaphore", 64, 0)
```
-or- to create a semaphore with a capacity of 64, and an initial count of 64:
```go
mysem := sem.Createsem("mySemaphore", 64, 64)
```

#### Createmutex: creates a mutex  
_func Createmutex(name string) *Mutex_
```go
mymutex := sem.Createmutex("myMutex")
```

Following a semaphore creation, there are a certain number of methods available to manipulate semaphores:

#### Reset
_func (s *Sem) Reset()_
```go
mysem.Reset()
```
-or- for a mutex

_func (m *Mutex) Reset()_
```go
mymutex.Reset()
```

This will flush the semaphore internal channel and resets its counter to its original value.

#### Signal -or- V (-or- Leave)
_func (s *Sem) Signal()_
```go
mysem.Signal()
```
-or-

_func (s *Sem) V()_
```go
mysem.V()
```
-or- for a mutex

_func (m *Mutex) Leave()_
```go
mymutex.Leave()
```

<b>Signal</b> and <b>V</b> accomplish the same thing which is to increase by 1 the level of availability of the resource. <b>Leave</b> is identical but reserved for <i>mutex</i>.

#### Wait -or- P (-or- Enter)
_func (s *Sem) Wait()_
```go
mysem.Wait()
```
-or-

_func (s *Sem) P()_
```go
mysem.P()
```
-or- for a mutex

_func (m *Mutex) Enter()_
```go
mymutex.Enter()
```
<b>Wait</b> and <b>P</b> accomplish the same thing which is to decrease by 1 the level of availability of the resource. <b>Enter</b> is identical but reserved for <i>mutex</i>. When the semaphore counter reaches 0, the resource is no longer available, until a Signal (-or- a V) call is made by another goroutine.

#### Notes:
- The <b>P</b> / <b>V</b> notation comes from <b>Edsger Dijkstra</b>, who introduced the concept of semaphores in 1963. The letters are from the Dutch words <b>Probeer</b> (try) and <b>Verhoog</b> (increment).

- The terms <b>Enter</b> and <b>Leave</b> for a mutex refer to ```Entering``` and ```Leaving``` critical sections in your code. A <b>critical section</b> is a region in your code that can be executed only by one goroutine at a time. Typically, you will need to define a critical section everytime you need to access a resource that can potentially be modified by multiple goroutines. Once in the critical section, a goroutine is guaranteed to have exclusive access to the shared resource.

