package sync

import (
	sem "github.com/crmathieu/gosem/pkg/semaphore"
	"fmt"
)

const MAXTHREADS = 100

type Sync struct {
	syncsem *sem.Sem
}

func NewSync(maxthreads int) *Sync {
	if maxthreads <= 0 || maxthreads > MAXTHREADS {
		maxthreads = MAXTHREADS
	}
	return &Sync{sem.Createsem("sync-sem", maxthreads, 0)}
}

func (s *Sync)Add(n int) {
	for i:=0; i<n; i++ {
		s.syncsem.Wait()
	}
}

func (s *Sync)Done(n int) {
	for i:=0; i<n; i++ {
		s.syncsem.Signal()
	}
}