package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const LastElement = "└───"
const DefaultElemtnt = "├───"

type ByName []os.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

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
	var ident, result string
	result = printDir(path, result, ident, printFiles)
	fmt.Fprintln(out, result)

	return nil
}
func getTab(index, lastIndex int) string {
	if index == lastIndex {
		return "\t"
	}

	return "│\t"
}
func fileSize(file os.FileInfo) string {
	if file.IsDir() {
		return ""
	} else if file.Size() <= 0 {
		return " (empty)"
	} else {
		return " (" + fmt.Sprint(file.Size()) + "b)"
	}
}

func getSortFiles(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(ByName(files))

	return files
}

func printDir(path string, result string, ident string, printFiles bool) string {
	sortFiles := getSortFiles(path)

	for index, f := range sortFiles {
		var line string
		if !printFiles && !f.IsDir() {
			continue
		}

		if index == len(sortFiles)-1 {
			line += ident + LastElement
		} else {
			line += ident + DefaultElemtnt
		}

		//fmt.Printf(line + f.Name() + fileSize(f) + "\n")
		result += line + f.Name() + fileSize(f) + "\n"

		if f.IsDir() {
			newPath := filepath.Join(path, f.Name())
			printDir(newPath, result, ident+getTab(index, len(sortFiles)-1), printFiles)
		}
	}

	return result
}
