package main

import (
	"github.com/mneumi/reading-crawler/engine"
	"github.com/mneumi/reading-crawler/task"
)

func main() {
	e := engine.New(10)

	e.Run([]task.Task{
		{
			URL:     "https://xxx.com",
			Handler: nil,
		},
	})
}
