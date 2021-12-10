package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// 作者模型
type Author struct {
	Mid         string           `json:"mid"`       // 作者id
	Name        string           `json:"name"`      // 作者名称
	Image       string           `json:"image"`     // 作者头像
	Description string           `json:"descption"` // 作者签名
	Gender      string           `json:"sex"`       // 作者性别
	Rank        int16            `json:"Rank"`
	Level       int16            `json:"level"`
	Num         int              `json:"num"`    // 视频数
	Videos      map[string]Video `json:"Videos"` // 视频列表
}

func NewAuthor(mid string) Author {
	return Author{Mid: mid, Videos: map[string]Video{}}
}

// 获取b站作者信息
// 每个主页查询结束会sleep 1s 防止频率太高被封禁
func (a *Author) GetInfo() error {
	// 获取作者信息
	err := a.author()
	if err != nil {
		fmt.Println("read author info fail:" + err.Error())
		return err
	}

	// 获取主页视频
	more := true
	for pn := 1; more; pn++ {
		more, err = a.homepage(pn, 30)
		if err != nil {
			fmt.Print("read home page fail, page num is " + strconv.Itoa(pn))
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// 获取主页内容 pn是页码 一页30个元素
// 返回结果 bool 为是否有下一页
// error 执行异常
// 当执行异常是 bool为false 直接中断操作
// 正常操作时 会填充author内信息
// 遍历完所有的主页内容, 用户的视频信息 以及视频信息就会构造完成
func (a *Author) homepage(pn int, ps int) (bool, error) {
	// eg. "https://api.bilibili.com/x/space/arc/search?mid=21869937&ps=40&order=pubdate"
	u := "https://api.bilibili.com/x/space/arc/search"
	params := url.Values{}
	params.Add("mid", a.Mid)
	params.Add("ps", strconv.Itoa(ps))
	params.Add("pn", strconv.Itoa(pn))
	params.Add("order", "pubdate")

	req, _ := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
	req.Header.Set("host", "api.bilibili.com")
	req.Header.Set("origin", "https://space.bilibili.com")
	req.Header.Set("refer", "https://space.bilibili.com/"+a.Mid)

	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("get videos of " + a.Mid + " fail: " + err.Error())
		return false, errors.New("request execute fail: " + err.Error())
	}

	if resp.StatusCode != 200 {
		fmt.Println("get videos of " + a.Mid + " fail: get statuscode(" + strconv.Itoa(resp.StatusCode) + ")")
		return false, errors.New("server status abnormal: get statuscode(" + strconv.Itoa(resp.StatusCode) + ")")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	type Homepage struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			List struct {
				Tlist map[string]struct {
					Tid   int    `json:"tid"`
					Name  string `json:"name"`
					Count int    `json:"count"`
				} `json:"tlist"`
				Vlist []Video `json:"vlist"`
			} `json:"list"`
			Page struct {
				Pn    int `json:"pn"`
				Ps    int `json:"ps"`
				Count int `json:"count"`
			} `json:"page"`
		} `json:"data"`
	}
	var h Homepage
	err = decoder.Decode(&h)
	if err != nil {
		fmt.Println("parse home page failed, error info: " + err.Error() +
			" pn: " + strconv.Itoa(pn) + " mid: " + a.Mid)
		return false, errors.New("parse json fail: " + err.Error())
	}

	if h.Code != 0 {
		fmt.Println("homepage message: [" + h.Msg + "] code:[" + strconv.Itoa(h.Code) + "]")
		return false, errors.New("homepage message: [" + h.Msg + "] code:[" + strconv.Itoa(h.Code) + "]")
	}

	if a.Num == 0 {
		a.Num = h.Data.Page.Count
	}

	//防止抓取的时候作者有上传 导致资源重复
	for _, video := range h.Data.List.Vlist {
		a.Videos[video.Bvid] = video
	}
	if a.Num-pn*ps > 0 {
		return true, nil
	}
	return false, nil
}

// 获取作者信息
func (a *Author) author() error {
	u := "https://api.bilibili.com/x/space/acc/info"
	params := url.Values{}
	params.Add("mid", a.Mid)

	req, _ := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
	req.Header.Set("host", "api.bilibili.com")
	req.Header.Set("origin", "https://space.bilibili.com")
	req.Header.Set("refer", "https://space.bilibili.com/"+a.Mid)

	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("get author info of  " + a.Mid + " fail: " + err.Error())
		return errors.New("request execute fail: " + err.Error())
	}

	if resp.StatusCode != 200 {
		fmt.Println("get author of " + a.Mid + " fail: get statuscode(" + strconv.Itoa(resp.StatusCode) + ")")
		return errors.New("server status abnormal: get statuscode(" + strconv.Itoa(resp.StatusCode) + ")")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	type AuthInfo struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			Name        string `json:"name"` // 姓名
			Gender      string `json:"sex"`  // 性别
			Image       string `json:"face"` // 头像
			Description string `json:"sign"` // 签名
			Rank        int16  `json:"Rank"`
			Level       int16  `json:"level"`
		} `json:"data"`
	}

	var ai AuthInfo
	err = decoder.Decode(&ai)
	if err != nil {
		fmt.Println("parse author info fail, error info: " + err.Error() + " mid: " + a.Mid)
		return errors.New("parse json fail: " + err.Error())
	}

	a.Name = ai.Data.Name
	a.Gender = ai.Data.Gender
	a.Image = ai.Data.Image
	a.Description = ai.Data.Description
	a.Level = ai.Data.Level
	a.Rank = ai.Data.Rank
	return nil
}
