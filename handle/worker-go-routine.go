package handle

import (
	"fmt"
	"time"
)

// TODO this file is using for manage go routine, not spawn undesirable  go routine like current behavior
// we can check https://gobyexample.com/worker-pools

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
