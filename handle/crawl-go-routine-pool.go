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
type WorkerRoutine struct {
	concurrentGoroutine int32
	jobQueue            []jobFn
	mu                  sync.Mutex
	limitAmountResource chan int
}

func NewWorkerRoutine(concurrentGoroutine int32) WorkerRoutine {
	if concurrentGoroutine < 1 {
		panic(errors.New("the concurrentGoroutine can be less than 1"))
	}
	limitAmountResource := make(chan int, concurrentGoroutine)
	var jobQueue []jobFn
	return WorkerRoutine{concurrentGoroutine: concurrentGoroutine,
		jobQueue: jobQueue, limitAmountResource: limitAmountResource}
}

// Execute TODO thanks chatgpt: https://chatgpt.com/, but idea is mine :)))
// TODO change the code later, use LinkedList to implement the queue, we do not need to copy
func (w *WorkerRoutine) Execute(job jobFn) {
	w.mu.Lock()

	w.jobQueue = append(w.jobQueue, job)
	jobExecutions := make([]jobFn, len(w.jobQueue))
	copy(jobExecutions, w.jobQueue)
	w.jobQueue = w.jobQueue[:0]

	w.mu.Unlock()

	for _, jobDetail := range jobExecutions {
		w.limitAmountResource <- 1
		go func(doJob jobFn) {
			defer func() { <-w.limitAmountResource }()
			doJob()
		}(jobDetail)
	}
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
