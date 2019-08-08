package main

import (
	sync "github.com/crmathieu/gosem/pkg/sync"
	"time"
	"fmt"
)

func main() {
	s := sync.NewSync(5)

	for i:=0; i<5; i++ {
		go func(a int) {
			defer s.Done()
			time.Sleep(time.Duration(a) * time.Second)
			fmt.Println("Finishing", a)
		}(i)
		s.Add()
	}
	s.Wait()
	fmt.Println("Main Finished...")
}