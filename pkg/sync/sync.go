package sync

import (
	sem "github.com/crmathieu/gosem/pkg/semaphore"
	"fmt"
)

type Sync struct {
	syncmutex *sem.Mutex
	csync   chan int
	cnt 	int
	in_use 	bool  
}

func NewSync() *Sync {
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
		fmt.Println("Invalid Done op. Make sure Each call to 'Done' has its corresponding call to 'Add'")
		return
	}
	s.cnt--
	left := s.cnt
	s.syncmutex.Leave()
	if left == 0 {
		s.csync <- 1
	} 
}

func (s *Sync)Wait() {
	<- s.csync
}