package queue

import (
  "github.com/crmathieu/gosem/semaphore"
)

type subQueue struct {
  Queue []int
  Head, Tail, Instance, Size int
  SemW *semaphore.Sem
  SemR *semaphore.Sem
}

type QU struct {
  consmap map[int]*subQueue
  cinstance int
  qmutex *semaphore.Mutex
  size   int
}

// Enter-----------------------------------------------------------------------
func (q *QU) Enter() {
  q.qmutex.Enter()
}

// Leave-----------------------------------------------------------------------
func (q *QU) Leave() {
  q.qmutex.Leave()
}

// InitQueue-------------------------------------------------------------------
func InitQueue(size int) *QU {
  return &QU{
    consmap : make(map[int]*subQueue),
    cinstance : 0,
    qmutex : semaphore.Cmutex("queueMutex"),
    size: size,
  }
}

// Next------------------------------------------------------------------------
func (sq *subQueue) Next() {
  sq.Tail = (sq.Tail + 1) % sq.Size
}

// Subscribe-------------------------------------------------------------------
func (q *QU) Subscribe() *subQueue {
  q.Enter()
  sq := subQueue{Size: q.size,
    Queue: make([]int, q.size),
    Instance: q.cinstance,
    Head: 0,
    Tail: 0,
    SemW: semaphore.Csem("fanbufW", q.size, q.size),
    SemR: semaphore.Csem("fanbufR", q.size, 0)}
  q.consmap[q.cinstance] = &sq
  q.cinstance++
  q.Leave()
  return &sq
}

// Close-----------------------------------------------------------------------
func (q *QU) Close(sq *subQueue) {
  q.Enter()
  delete(q.consmap, sq.Instance)
  q.Leave()
}


// Read------------------------------------------------------------------------
func (sq *subQueue) Read() (item int, index int) {
  sq.SemR.Wait()
  item, index = sq.Queue[sq.Tail], sq.Tail
  sq.Tail = (sq.Tail + 1) % sq.Size
  sq.SemW.Signal()
  return
}

// GetIndex--------------------------------------------------------------------
func (sq *subQueue) GetIndex() int {
  return sq.Tail
}

// Fanout----------------------------------------------------------------------
func (q *QU) Fanout(item int) {
  var max = len(q.consmap)
  for k := 0; k < max; k++ {
    q.consmap[k].SemW.Wait()
    q.Enter()
    q.consmap[k].Queue[q.consmap[k].Head] = item
    q.consmap[k].Head = (q.consmap[k].Head + 1) % q.consmap[k].Size
    q.Leave()
    max = len(q.consmap)
    q.consmap[k].SemR.Signal()
  }
}
