package hello

import (
	"fmt"
	"os"
)

func main()  {
	fmt.Println("arg size:",len(os.Args))
	fmt.Println("hello world")
}