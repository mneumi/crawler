package task

type Task struct {
	URL     string
	Handler func(string) TaskHandleResult
}

type TaskHandleResult struct {
	Info  interface{}
	Tasks []Task
}
