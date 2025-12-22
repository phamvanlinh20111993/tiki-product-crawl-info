package handle

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type jobFn func()

// TODO this file is using for manage go routine, not spawn undesirable  go routine like current behavior
// we can check https://gobyexample.com/worker-pools
type WorkerPool struct {
	concurrentGoroutine int32
	jobQueue            []jobFn
	mu                  sync.Mutex
	limitAmountResource chan int
}

func NewWorkerPool(concurrentGoroutine int32) WorkerPool {
	if concurrentGoroutine < 1 {
		panic(errors.New("workers must be >= 1"))
	}
	limitAmountResource := make(chan int, concurrentGoroutine)
	var jobQueue []jobFn
	return WorkerPool{concurrentGoroutine: concurrentGoroutine,
		jobQueue: jobQueue, limitAmountResource: limitAmountResource}
}

// Execute TODO thanks chatgpt: https://chatgpt.com/, but idea is mine :)))
// TODO change the code later, use LinkedList to implement the queue, we do not need to copy
func (wp *WorkerPool) Execute(job jobFn) {
	wp.mu.Lock()

	wp.jobQueue = append(wp.jobQueue, job)
	jobExecutions := make([]jobFn, len(wp.jobQueue))
	copy(jobExecutions, wp.jobQueue)
	wp.jobQueue = wp.jobQueue[:0]

	wp.mu.Unlock()

	for _, jobDetail := range jobExecutions {
		wp.limitAmountResource <- 1
		go func(doJob jobFn) {
			defer func() { <-wp.limitAmountResource }()
			doJob()
		}(jobDetail)
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.limitAmountResource)
}

// TODO mark as remove later
func worker(jobId int, ch chan int) {
	fmt.Println("JobId ", jobId, "start")
	for i := 0; i < 5; i++ {
		fmt.Println("JobId ", jobId, " wake up and then go to sleep")
		time.Sleep(2 * time.Second)
	}
	fmt.Println("JobId ", jobId, "end")
	<-ch
}

func jobs(ch chan int) {
	for i := 0; i < 100; i++ {
		go worker(i, ch)
		ch <- i
	}
	close(ch)
}

func Example() {
	ch := make(chan int, 4)
	jobs(ch)
}
