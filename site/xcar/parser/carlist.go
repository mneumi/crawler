package parser

import (
	"regexp"

	"github.com/mneumi/reading-crawler/task"
)

const host = "http://newcar.xcar.com.cn"

var carModelRe = regexp.MustCompile(`<a href="(/\d+/)" target="_blank" class="list_img">`)
var carListRe = regexp.MustCompile(`<a href="(//newcar.xcar.com.cn/car/[\d+-]+\d+/)"`)

func ParseCarList(contents []byte) *task.TaskHandleResult {
	result := &task.TaskHandleResult{}

	matches := carModelRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Tasks = append(result.Tasks, task.Task{
			URL: host + string(m[1]),
			Handler: func(contents []byte) *task.TaskHandleResult {
				return ParseCarModel(contents, string(m[1]))
			},
		})
	}

	matches = carListRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Tasks = append(result.Tasks, task.Task{
			URL:     "http:" + string(m[1]),
			Handler: ParseCarList,
		})
	}

	return result
}
