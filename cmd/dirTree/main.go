package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	path = "./testdata"
	printDir(path)
	//files, err := ioutil.ReadDir(path)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, f := range files {
	//
	//	if f.IsDir() {
	//		filesInDir, err := ioutil.ReadDir(path + "/" +f.Name())
	//
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		for _, ff := range filesInDir {
	//			fmt.Printf("├───%s\n" ,ff.Name())
	//		}
	//
	//		fmt.Printf("├───%s\n" ,f.Name())
	//	}
	//}

	return nil
}

func printDir(path string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	count := strings.Count(path, "/")

	for _, f := range files {
		if f.IsDir() {
			twoPath := path + "/" + f.Name()
			printDir(twoPath)
		}
		if 0 < count {
			//tab := strings.Repeat("\n", count)

			fmt.Printf("\t├───%s\n", f.Name())
		} else {
			fmt.Printf("├───%s\n", f.Name())
		}
	}
}
