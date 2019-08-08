package sync

import (
	sem "github.com/crmathieu/gosem/pkg/semaphore"
	"fmt"
)

const MAXTHREADS = 100

type Sync struct {
	syncmutex *sem.Mutex
	csync   chan int
	cnt 	int
	in_use 	bool  
}

func NewSync(maxthreads int) *Sync {
	if maxthreads <= 0 || maxthreads > MAXTHREADS {
		maxthreads = MAXTHREADS
	}
	return &Sync{sem.Createmutex("sync-mutex"), make(chan int, 1), 0, false}
}

func (s *Sync)Add() {
	s.syncmutex.Enter()
	if s.in_use == false {
		s.in_use = true
	}
	s.cnt++
	s.syncmutex.Leave()
}

func (s *Sync)Done() {
	s.syncmutex.Enter()
	if s.in_use == false || s.cnt < 1 {
		fmt.Println("Invalid Done op")
		return
	}
	if s.cnt == 1 {
		s.csync <- 1
	} 
	s.cnt--
	s.syncmutex.Leave()
}

func (s *Sync)Wait() {
	<- s.csync
}