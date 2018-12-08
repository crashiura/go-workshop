package main

import (
	"flag"
	"fmt"
	"os"
)

var greeterType = flag.String("greeter", "suffix", "type of greeter")

type Greeter interface {
	Greet(name string)
}

type PrefixGreeter struct {
}

func (p PrefixGreeter) Greet(name string) {
	fmt.Println("hello", name)
}

type SuffixGreeter struct {
}

func (s SuffixGreeter) Greet(name string) {
	fmt.Println(name, "hello")
}

func main() {
	//
	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: %s <type-of-greeter> \n", os.Args[0])
	//	os.Exit(1)
	//}
	flag.Parse()

	var g Greeter
	switch *greeterType {
	case "prefix":
		g = PrefixGreeter{}
	case "suffix":
		g = SuffixGreeter{}
	default:
		fmt.Println("exit")
		os.Exit(1)
	}
	g.Greet("world")
}
