package engine

import (
	"log"
	"time"

	"github.com/mneumi/reading-crawler/scheduler"
	"github.com/mneumi/reading-crawler/task"
	"github.com/mneumi/reading-crawler/worker"
)

type Engine struct {
	in          chan *task.Task
	out         chan *task.TaskHandleResult
	workerCount int
	workers     []*worker.Worker
	scheduler   *scheduler.Scheduler
}

func New(workerCount int) *Engine {
	in := make(chan *task.Task, 10)
	out := make(chan *task.TaskHandleResult, 10)

	e := &Engine{
		in:          in,
		out:         out,
		workerCount: workerCount,
	}

	e.bindScheduler()
	e.bindWorkers()

	return e
}

func (e *Engine) Run(ts []task.Task) {
	for _, worker := range e.workers {
		worker.Start()
	}

	for _, t := range ts {
		e.scheduler.Submit(&t)
	}

	e.processResult()
}

func (e *Engine) bindScheduler() {
	e.scheduler = scheduler.New(e.in)
}

func (e *Engine) bindWorkers() {
	for i := 0; i < e.workerCount; i++ {
		e.workers = append(e.workers, worker.New(i+1, e.in, e.out))
	}
}

func (e *Engine) processResult() {
	for {
		select {
		case result := <-e.out:
			// process result here
			if len(result.Info) > 0 {
				log.Printf("\n\n%+v\n\n", result.Info)
			}

			for _, t := range result.Tasks {
				tCopy := t
				e.scheduler.Submit(&tCopy)
			}
		// 如果 2 分钟内，out 都没有收到新的 result，证明所有任务结束，程序退出
		case <-time.NewTicker(2 * time.Minute).C:
			log.Println("all done")
			return
		}
	}
}
