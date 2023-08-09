package main

import (
	"fmt"
	"sync"
)

func usoMutex() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	contador := 0
	gorotinas := 10
	incrementador := 1000

	for i := 0; i < gorotinas; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementador; j++ {
				mu.Lock()
				contador++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Println("Valor final do contador", contador)
}
