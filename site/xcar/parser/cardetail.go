package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mneumi/reading-crawler/site/xcar/model"
	"github.com/mneumi/reading-crawler/task"
)

var nameRe = regexp.MustCompile(`<title>【(.*)】报价_图片_参数-爱卡汽车.*</title>`)
var imageRe = regexp.MustCompile(`<img class="color_car_img_new" src="([^"]+)"`)
var priceReTmpl = `<a href="/%s/baojia/".*>(\d+\.\d+)</a>`

func ParseCarDetail(contents []byte, url string, brandId string) *task.TaskHandleResult {
	id := strings.ReplaceAll(url, "/", "")

	car := model.CarDetail{
		BrandId:  brandId,
		Name:     extractString(contents, nameRe),
		ImageURL: "http:" + extractString(contents, imageRe),
	}

	priceRe, err := regexp.Compile(fmt.Sprintf(priceReTmpl, id))

	if err == nil {
		car.Price = extractFloat(contents, priceRe)
	}

	result := &task.TaskHandleResult{
		Info: []interface{}{car},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

func extractFloat(contents []byte, re *regexp.Regexp) float64 {
	f, err := strconv.ParseFloat(extractString(contents, re), 64)

	if err != nil {
		return 0
	}

	return f
}
