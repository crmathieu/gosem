package main

import (
	sync "github.com/crmathieu/pkg/sync"
	"time"
)

func main() {
	s := sync.NewSync(5)

	for i:=0; i<5; i++ {
		go func() {
			time.Sleep(1 * time.Duration.Second)
		}()
	}
}