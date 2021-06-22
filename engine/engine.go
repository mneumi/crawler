package engine

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/mneumi/reading-crawler/persist"
	"github.com/mneumi/reading-crawler/scheduler"
	"github.com/mneumi/reading-crawler/site/xcar/model"
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
		// ts 也进行去重的原因是把源URL也放入去重表中
		if isDuplicate(t.URL) {
			log.Printf("Duplicate request: %s", t.URL)
			continue
		}

		tCopy := t
		e.scheduler.Submit(&tCopy)
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
			if result.Info != nil {
				if item, ok := result.Info.(model.CarModel); ok {
					persist.DB.Table("car_models").Create(&item)
				}
				if item, ok := result.Info.(model.CarDetail); ok {
					persist.DB.Table("car_details").Create(&item)
				}
			}

			for _, t := range result.Tasks {
				if isDuplicate(t.URL) {
					log.Printf("Duplicate request: %s", t.URL)
					continue
				}

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

func isDuplicate(url string) bool {
	_, err := persist.RDB.Get(url).Result()

	// 如果 err != redis.Nil，则证明库里有这个值，即重复了，直接返回
	if err != redis.Nil {
		return true
	}

	// 不重复，则将键加入 redis
	err = persist.RDB.Set(url, 1, -1).Err()

	if err != nil {
		log.Printf("set score failed, err:%v\n", err)
	}

	return false
}
