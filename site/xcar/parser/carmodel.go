package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mneumi/reading-crawler/site/xcar/model"
	"github.com/mneumi/reading-crawler/task"
)

var logoRe = regexp.MustCompile(`<meta itemprop="image" content="//img3.xcarimg.com/PicLib/logo/([^_]+)_40.jpg" />`)
var brandRe = regexp.MustCompile(`<span class="lt_f1">([^<]+)</span><h1>([^<]+)<a`)
var carDetailRe = regexp.MustCompile(`<a href="(/m\d+/)" target="_blank"`)

func ParseCarModel(contents []byte, url string) *task.TaskHandleResult {
	logoMatch := logoRe.FindSubmatch(contents)
	logoURL := fmt.Sprintf("http://img3.xcarimg.com/PicLib/logo/%s_160.jpg", logoMatch[1])

	productMatch := brandRe.FindSubmatch(contents)
	brandId := strings.ReplaceAll(url, "/", "")
	brandName := productMatch[1]
	brandModel := productMatch[2]

	matches := carDetailRe.FindAllSubmatch(contents, -1)

	result := &task.TaskHandleResult{
		Info: []interface{}{
			model.CarModel{
				LogoURL:    logoURL,
				BrandName:  string(brandName),
				BrandModel: string(brandModel),
			},
		},
	}

	for _, m := range matches {
		url := string(m[1])

		result.Tasks = append(result.Tasks, task.Task{
			URL: host + url,
			Handler: func(contents []byte) *task.TaskHandleResult {
				return ParseCarDetail(contents, url, brandId)
			},
		})
	}

	return result
}
