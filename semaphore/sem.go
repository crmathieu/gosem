package semaphore

type Sem struct {
  name string
  sch chan int
  ic int
  valid bool
<<<<<<< HEAD
}

type Mutex struct {
  s *Sem
}

// Cmutex----------------------------------------------------------------------
//         creates a mutex
// ----------------------------------------------------------------------------
func Cmutex(name string) *Mutex {
  m := &Mutex{
    s: Csem(name, 1, 1),
  }  
  return m
=======
}

type Mutex struct {
  s *Sem
>>>>>>> 1ef29cbbc4cb55336c700fcf2b5a69b3ae322f1b
}

// Csem------------------------------------------------------------------------
// creates a semaphore with an initial count (ic) >= 0. The initial count must
// lower or equal to the semaphore capacity
// ----------------------------------------------------------------------------
func Csem(name string, capacity int, ic int) *Sem {

  if capacity <= 0 {
    capacity = 1
  }
  if ic < 0 {
    ic = 0
  }
  if ic > capacity {
    ic = capacity
  }

  s := &Sem{name: name, 
            sch: make(chan int, capacity), 
            ic: ic, 
            valid: true,
          }
  for i := 0; i < ic; i++ {
    s.sch <- 1
  }
  return s
}

// Close-----------------------------------------------------------------------
<<<<<<< HEAD
//         Close semaphore. This will release resources waiting for it
=======
// Closes a semaphore. This will release resources waiting for it
>>>>>>> 1ef29cbbc4cb55336c700fcf2b5a69b3ae322f1b
// ----------------------------------------------------------------------------
func (s *Sem) Close() {
  if s.valid {
    s.valid = false
    close(s.sch)
  }
}


// Reset-----------------------------------------------------------------------
// resets the semaphore.
// ----------------------------------------------------------------------------
func (s *Sem) Reset() {
  _flushChannel(s)
  for i := 0; i < s.ic; i++ {
    s.sch <- 1
  }
}

// _flushChannel---------------------------------------------------------------
func _flushChannel(s *Sem) {
  for {
    select {
      case <- s.sch:
      default:
        return
    }
  }
}

<<<<<<< HEAD
// Signal----------------------------------------------------------------------
//         releases the semaphore. V & Signal as synonymous
// ----------------------------------------------------------------------------
=======
// Signal-----------------------------------------------------------------------
// releases the semaphore. V & Signal as synonymous. Leave is the mutex equival
// -----------------------------------------------------------------------------
>>>>>>> 1ef29cbbc4cb55336c700fcf2b5a69b3ae322f1b
func (s *Sem) V() bool {
  s.sch <- 1
  return s.valid
}

func (s *Sem) Signal() bool {
  s.sch <- 1
  return s.valid
<<<<<<< HEAD
}

func (m *Mutex) Leave() bool {
    m.s.sch <- 1
  return m.s.valid
=======
>>>>>>> 1ef29cbbc4cb55336c700fcf2b5a69b3ae322f1b
}


// Wait------------------------------------------------------------------------
// acquires the semaphore. P & Wait are synonymous. Enter is the mutex equival
// ----------------------------------------------------------------------------
func (s *Sem) P() bool {
  <- s.sch
  return s.valid
}

func (s *Sem) Wait() bool {
  <- s.sch
  return s.valid
}

<<<<<<< HEAD
func (m *Mutex) Enter() bool {
   <- m.s.sch
  return m.s.valid
=======

// Cmutex----------------------------------------------------------------------
// creates a mutex. A mutex is a semaphore with capacity and initial count = 1
// ----------------------------------------------------------------------------
func Cmutex(name string) *Mutex {
  m := &Mutex{
    s: Csem(name, 1, 1),
  }  
  return m
}

// Enter-----------------------------------------------------------------------
// enters a critical section. Equivalent to P or Wait
//-----------------------------------------------------------------------------
func (m *Mutex) Enter() bool {
   <- m.s.sch
  return m.s.valid
}

// Leave-----------------------------------------------------------------------
// leaves a critical section. Equivalent to V or Signal
//-----------------------------------------------------------------------------
func (m *Mutex) Leave() bool {
    m.s.sch <- 1
  return m.s.valid
}

// Close-----------------------------------------------------------------------
// Closes a mutex. This will release resources waiting for it
// ----------------------------------------------------------------------------
func (m *Mutex) Close() {
  if m.s.valid {
    m.s.valid = false
    close(m.s.sch)
  }
>>>>>>> 1ef29cbbc4cb55336c700fcf2b5a69b3ae322f1b
}


// Reset-----------------------------------------------------------------------
// resets a mutex.
// ----------------------------------------------------------------------------
func (m *Mutex) Reset() {
  _flushChannel(m.s)
  m.s.sch <- 1
}
