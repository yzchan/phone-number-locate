// Package work 来自于《Go in Action》7.3章节并发模式中的work包
// 本包使用无缓冲的通道来创建一个goroutine池，这些goroutine执行并控制一组工作，让其并发执行
package work

import "sync"

type Tasker interface {
	Task()
}

type Pool struct {
	tasks chan Tasker
	wg    sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{
		tasks: make(chan Tasker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.tasks {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

func (p *Pool) Run(w Tasker) {
	p.tasks <- w
}

func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
