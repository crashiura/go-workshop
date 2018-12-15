package pool

import "sync"

type Pool struct {
	// Parallelism is a maximum number of spawned goroutines.
	// If Parallelism is zero then the default value of one used.
	Parallelism int

	initOnce sync.Once
	stopOnce sync.Once

	sem chan struct{}
	wrk chan func()
}

func (p *Pool) init() {
	p.initOnce.Do(func() {
		n := p.Parallelism
		if n == 0 {
			n = 1
		}
		p.sem = make(chan struct{}, n)
		p.wrk = make(chan func())
	})
}

func (p *Pool) stop() {
	p.stopOnce.Do(func() {
		close(p.wrk)
		for i := 0; i < cap(p.sem); i++ {
			p.sem <- struct{}{}
		}
	})
}

func (p *Pool) Schedule(task func()) {
	p.init()

	select {
	case p.wrk <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

// Close stops the scheduling of tasks and waits for all workers are done.
func (p *Pool) Close() {
	p.init()
	p.stop()
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()
	for task != nil {
		task()
		task = <-p.wrk
	}
}
