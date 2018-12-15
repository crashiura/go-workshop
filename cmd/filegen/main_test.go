package main

import (
	"io/ioutil"
	"strconv"
	"testing"
)

type stubFileWriter struct {
	writes [][]byte
}

func (s *stubFileWriter) Write(p []byte) (int, error) {
	cp := make([]byte, len(p))
	copy(cp, p)
	s.writes = append(s.writes, cp)
	return len(p), nil
}

type stubRandSource struct {
	seed    int64
	curr    int64
	numbers []int64
}

func (s *stubRandSource) Int63() (r int64) {
	r = s.curr
	s.numbers = append(s.numbers, s.curr)
	s.curr++
	return
}

func (s *stubRandSource) Seed(v int64) {
	s.seed = v
}

func TestWriteRandNumber(t *testing.T) {
	var s stubFileWriter
	var r stubRandSource

	rnd := randNumberWriter{
		Source: &r,
	}
	err := rnd.WriteTo(&s)
	if err != nil {
		t.Fatal(err)
	}
	if act, exp := len(s.writes), 1; act != exp {
		t.Fatalf(
			"unexpected number of Write() calls: %d; want %d",
			act, exp,
		)
	}
	p := s.writes[0]
	n := len(p)
	if n < 1 || p[n-1] != '\n' {
		t.Fatalf(
			"malformed bytes written: %+q", p,
		)
	}
	num := p[:n-1]
	act, err := strconv.ParseInt(string(num), 10, 64)
	if err != nil {
		t.Fatalf("malformed number written: %v", err)
	}
	if len(r.numbers) < 1 {
		t.Fatalf("Source() was not called")
	}
	if exp := r.numbers[0]; act != exp {
		t.Fatalf("unexpected number written: %d; want %d", act, exp)
	}
}

type noopSource struct{}

func (noopSource) Seed(int64)   {}
func (noopSource) Int63() int64 { return 0 }

func BenchmarkWriteRandNumber(b *testing.B) {
	rnd := randNumberWriter{
		Source: noopSource{},
	}
	for i := 0; i < b.N; i++ {
		rnd.WriteTo(ioutil.Discard)
	}
}
