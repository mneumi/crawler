package worker

import (
	"log"

	"github.com/mneumi/reading-crawler/fetcher"
	"github.com/mneumi/reading-crawler/limiter"
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
			// 限流
			<-limiter.GetLimiter()

			t := <-w.in

			// fetch & parse here
			result, err := w.do(t)

			if err != nil {
				log.Printf("ERR: %v", err)
			}

			if result == nil {
				continue
			}

			w.out <- result
		}
	}()
}

func (w *Worker) do(t *task.Task) (*task.TaskHandleResult, error) {
	log.Printf("#%d work start ... %v", w.ID, t.URL)

	body, err := fetcher.Fetch(t.URL)

	if err != nil {
		log.Printf("fetch error url: %s, err: %v", t.URL, err)
		return &task.TaskHandleResult{}, err
	}

	if t.Handler != nil {
		return t.Handler(body), nil
	}

	return nil, nil
}
