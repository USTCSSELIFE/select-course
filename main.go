package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	requestUrl  string = "http://yjs.ustc.edu.cn/course/query.asp?mode=dept"
	notifyToken string = ""
)

var courseNumbers = []string{"EIEN6017P"}

func main() {
	u := launcher.New().NoSandbox(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustIncognito().MustPage(requestUrl).MustWindowFullscreen()
	frame := page.MustWaitLoad().MustElement("#query").MustFrame()
	frame.MustElement("#year1").MustSelect("2021-2022")
	frame.MustElement("#radio1").MustClick()
	frame.MustElement("#kkdept").MustSelect("A14-软件学院苏州")
	frame.MustElementX(`//input[@type="submit"]`).MustClick()
	time.Sleep(time.Second * 2)

	for _, cNo := range courseNumbers {
		courses := frame.MustElementsX(`//tr[@class="bt06" and contains(., "` + cNo + `")]`)
		for _, course := range courses {
			courseName := course.MustElementsX(`.//a[@title="详细信息"]`)[0].MustText()
			selectNumStr := course.MustElementsX(`.//a[@title="详细信息"]`)[1].MustText()
			fullNumStr := course.MustElementsX(`.//td[@class="bt06"]`)[6].MustText()
			courseTeacher := course.MustElementsX(`.//td[@class="bt06"]`)[2].MustText()
			courseTime := course.MustElementsX(`.//td[@class="bt06"]`)[3].MustText()
			selectNum, _ := strconv.Atoi(selectNumStr)
			fullNum, _ := strconv.Atoi(fullNumStr)
			if selectNum < fullNum {
				text := fmt.Sprintf("%s %s 周的%s有 %d 个空位。", courseTeacher, courseTime, courseName, fullNum-selectNum)
				notify(text)
			}
		}
	}
}

func notify(text string) {
	notifyUrl := fmt.Sprintf("https://api.day.app/%s/选课空位通知/%s", notifyToken, url.QueryEscape(text))
	_, _ = http.Get(notifyUrl)
}
