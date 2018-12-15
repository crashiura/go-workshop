package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

var (
	n int // number for files
	m int // number of entries per file
)

var (
	max = len(strconv.AppendUint(nil, math.MaxInt64, 10))
)

var fileNames [100]string

func init() {
	for i := range fileNames {
		fileNames[i] = fmt.Sprintf("%d.gen", i)
	}
}

func main() {
	flag.IntVar(&n, "n", 1, "number of files to generate")
	flag.IntVar(&m, "m", 5, "number of entries to generate file")
	flag.Parse()

	results := make(chan error)

	for i := 0; i < n; i++ {
		go fileWriter(results, i, m)
	}

}

func fileWriter(res chan<- error, i int, m int) {
	var err error
	file, err := os.Create(fileNames[i])

	//if err != nil {
	//	log.Fatal(err)
	//	res <-err
	//	return
	//}

	defer func() {
		if err != nil {
			log.Fatal(err)

			res <- err
		}
	}()

	for j := 0; j < m; j++ {
		bts := genereateNumberButes()

		_, err := file.Write(bts)

		if err != nil {
			log.Fatal("write number to file error: %v", err)
			res <- err
			return
		}
	}
}

func genereateNumberButes() []byte {
	v := rand.Uint64()
	bts := make([]byte, 0, max+1)
	bts = strconv.AppendUint(nil, v, 10)
	bts = append(bts, '\n')

	return bts
}
