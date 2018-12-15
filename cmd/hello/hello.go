package main

import (
	"flag"
	"fmt"
	"os"
)

var greeterType = flag.String("greeter", "suffix", "type of the greeter")

// Greeter greets someone by name.
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
	flag.Parse()

	var g Greeter
	switch *greeterType {
	case "prefix":
		g = PrefixGreeter{}
	case "suffix":
		g = SuffixGreeter{}
	default:
		fmt.Println("unexpected greeter:", os.Args[1])
		os.Exit(1)
	}

	g.Greet("world")
}
