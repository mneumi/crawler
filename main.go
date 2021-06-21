package main

import (
	"github.com/mneumi/reading-crawler/engine"
	"github.com/mneumi/reading-crawler/site/xcar/parser"
	"github.com/mneumi/reading-crawler/task"
)

func main() {
	e := engine.New(10)

	e.Run([]task.Task{
		{
			URL:     "https://newcar.xcar.com.cn/",
			Handler: parser.ParseCarList,
			// URL: "https://newcar.xcar.com.cn/3428/",
			// Handler: func(b []byte) *task.TaskHandleResult {
			// 	return parser.ParseCarModel(b, "")
			// },
		},
	})
}
