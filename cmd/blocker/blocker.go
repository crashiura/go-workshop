package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	go monitor()
	go cpueater()

	var block chan struct{}
	<-block
}

func cpueater() {
	for now := range time.Tick(time.Second * 5) {
		fmt.Println("eater started", now)
		for i := 0; i < 1e10; i++ {
			if i%1e9 == 0 {
				//fmt.Println("eater cooperates")
				runtime.Gosched()
			}
		}
		fmt.Println("eater done")
	}
}

func monitor() {
	const duration = 500 * time.Millisecond

	lastTick := time.Now()
	for now := range time.Tick(duration) {
		fmt.Println("time", now)
		waitDuration := now.Sub(lastTick)
		if waitDuration > 2*duration {
			fmt.Println(
				"WARNING: cpueater detected! Overdelayed to",
				waitDuration-duration,
			)
		}
		lastTick = now
	}
}
