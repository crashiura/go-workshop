package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gobwas/bd/pool"
)

var (
	limit = flag.Int("limit", 5, "max number of goroutines to use")
	word  = flag.String("word", "go", "word to search for")
)

func main() {
	flag.Parse()

	p := pool.Pool{
		Parallelism: *limit,
	}
	for _, s := range flag.Args() {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			log.Printf("skipping invalid addr: %q", s)
			continue
		}
		addr := s
		p.Schedule(func() {
			rc, err := fetch(addr)
			if err != nil {
				log.Fatal("fetch error: %v", err)
			}
			defer rc.Close()

			n, err := scan(rc, *word)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			log.Printf(
				"found %q at %s %d times",
				*word, addr, n,
			)
		})
	}
	p.Close()
}

func scan(r io.Reader, word string) (count int, err error) {
	var (
		w   = []byte(word)
		buf = make([]byte, 4096)
	)
	for {
		var n int
		n, err = r.Read(buf)
		data := buf[:n]
		for {
			i := bytes.Index(data, w)
			if i == -1 {
				break
			}
			count++
			data = data[i+len(w):]
		}
		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			return
		}
		// FIXME: handle split case.
	}
	return
}

func fetch(addr string) (io.ReadCloser, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
