package main

import (
	"fmt"
	"runtime"
	"time"
)

type ball struct{}

func main() {
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(1))
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	var b ball
	table := make(chan ball, 2)

	var done [2]chan struct{}
	for i, name := range []string{"Petr", "Ivan"} {
		done[i] = make(chan struct{})
		go player(table, done[i], name)
	}

	table <- b

	time.Sleep(time.Second)

	fmt.Println("time to stop the game")
	<-table
	close(table)

	fmt.Println("waiting for the players...")
	for _, done := range done {
		<-done
		fmt.Println("player gone")
	}
}

func player(table chan ball, done chan struct{}, name string) {
	defer close(done)
	for b := range table {
		fmt.Println(name, "YAY! Got the ball")
		time.Sleep(100 * time.Millisecond)
		table <- b
	}
}
