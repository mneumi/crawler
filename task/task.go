package task

type Task struct {
	URL     string
	Handler func([]byte) *TaskHandleResult
}

type TaskHandleResult struct {
	Info  []interface{}
	Tasks []Task
}
