package common

import (
	"fmt"
	"sync"
	"time"
)

// Start Start
func Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		quit := make(chan int)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ticker.C:
					fmt.Println("ticker .")
				case <-quit:
					fmt.Println("work well .")
					ticker.Stop()
					return
				}
			}
		}()
		time.Sleep(1)
		quit <- 1
		wg.Wait()
	}()
}
