package main

import (
	"fmt"
)

type ball struct {

}

func main()  {
	var b ball

	table := make(chan ball)

	go player("Ivan",table)
	go player("Peter",table)

	table <- b

	//time.Sleep(time.Second)
	fmt.Println("close game")

	<-table
	close(table)
}

func player(name string, table chan ball)  {
	for  b := range table  {
		fmt.Println("YAY! Got the ball", name)
		table <-b
		//time.Sleep(100 * time.Microsecond)
	}
}