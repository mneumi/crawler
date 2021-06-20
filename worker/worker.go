package worker

import (
	"log"
	"math/rand"
	"time"

	"github.com/mneumi/reading-crawler/task"
)

type Worker struct {
	ID  int
	in  chan *task.Task
	out chan *task.TaskHandleResult
}

func New(id int, in chan *task.Task, out chan *task.TaskHandleResult) *Worker {
	return &Worker{
		ID:  id,
		in:  in,
		out: out,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			t := <-w.in

			// fetch & parse here
			// TOOD
			log.Printf("#%d work start ... %v", w.ID, t.URL)
			log.Printf("#%d work done ...", w.ID)

			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(100)
			println(n)

			if n < 2 {
				w.out <- &task.TaskHandleResult{
					Info: "Hello World",
				}
			} else {
				w.out <- &task.TaskHandleResult{
					Info: "Hello World",
					Tasks: []task.Task{
						{
							URL:     "http://abcdefg.com",
							Handler: nil,
						},
					},
				}
			}
		}
	}()
}
