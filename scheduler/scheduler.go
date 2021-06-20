package scheduler

import "github.com/mneumi/reading-crawler/task"

type Scheduler struct {
	In chan *task.Task
}

func New(in chan *task.Task) *Scheduler {
	return &Scheduler{
		In: in,
	}
}

func (s *Scheduler) Submit(t *task.Task) {
	go func() {
		s.In <- t
	}()
}
