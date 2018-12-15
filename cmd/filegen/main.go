package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

var (
	n int // Number of files.
	m int // Number of entries per file.
)

var (
	max = len(strconv.AppendUint(nil, math.MaxUint64, 10))
)

var fileNames [100]string

func init() {
	for i := range fileNames {
		fileNames[i] = fmt.Sprintf("%d.gen", i)
	}
}

func main() {
	flag.IntVar(&n, "n", 1, "number of files to generate")
	flag.IntVar(&m, "m", 5, "number of entries to generate per file")
	flag.Parse()
	if n > 100 {
		log.Fatalf("maximum number of files reached")
	}

	results := make(chan error)

	for i := 0; i < n; i++ {
		go fileWriter(results, i, m)
	}
	for i := 0; i < n; i++ {
		err := <-results
		if err == nil {
			continue
		}
		createErr, ok := err.(CreateFileError)
		if ok {
			log.Fatalf(
				"create file #%d error: %v",
				createErr.Index, createErr.Reason,
			)
		}
		if err == ErrWriteRandomNumber {
			log.Fatalf(
				"write random number to a file error",
			)
		}
		panic("unexpected error type")
	}
}

type CreateFileError struct {
	Reason error
	Index  int
}

func (e CreateFileError) Error() string {
	return "create file error: " + e.Reason.Error()
}

var (
	ErrWriteRandomNumber = errors.New("write random number error")
)

func fileWriter(res chan<- error, i, m int) {
	var err error
	defer func() {
		res <- err
	}()

	file, err := os.Create(fileNames[i])
	if err != nil {
		err = CreateFileError{
			Reason: err,
			Index:  i,
		}
		return
	}
	defer file.Close()

	rnd := randNumberWriter{
		Source: rand.NewSource(1),
	}

	buf := bufio.NewWriter(file)
	for j := 0; j < m; j++ {
		if err = rnd.WriteTo(buf); err != nil {
			err = ErrWriteRandomNumber
			return
		}
	}
	err = buf.Flush()
}

type randNumberWriter struct {
	Source rand.Source
}

var pool = sync.Pool{
	New: func() interface{} {
		p := make([]byte, 0, max+1) // max is for uint64; 1 is for '\n'.
		return &p
	},
}

func (r randNumberWriter) WriteTo(w io.Writer) error {
	v := r.Source.Int63()

	bts := pool.Get().(*[]byte)
	defer pool.Put(bts)

	p := *bts
	p = strconv.AppendInt(p, v, 10)
	p = append(p, '\n')
	_, err := w.Write(p)

	return err
}
