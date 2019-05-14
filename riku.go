package main

import (
	"encoding/json"
	"fmt"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// 二维码显示在终端上
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		logs.Error(err)
		return
	}

	session.SetAfterLogin(func() error{
		qi := session.Cm.GetContactsByName("琪琪快乐买货宝")
		logs.Informational("qi: %+v", qi)
		if len(qi) >0 {
			logs.Informational("qi0: %+v", qi[0])
			qq := qi[0]
			ticker := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-ticker.C:
					logs.Info("定时发消息")
					//logs.Informational("m name %s", session.Bot.UserName)
					//logs.Informational("g name %s", qq.UserName)
					m := fmt.Sprintf("%s\n%s",
						time.Now().Format("2006/1/2 15:04:05"),
						Worker())
					session.SendText(m, session.Bot.UserName, qq.UserName)
				}
			}
		} else {
			return nil
		}
		return nil
	})

	// 登录并接收消息
	if err := session.LoginAndServe(false); err != nil {
		logs.Error("session exit, %s", err)
	}
}

type Stock struct {
	Key    string `json:"key"`
	SkuId  string `json:"SKUID"`
	Value  string `json:"value"`
	Stock  string `json:"Stock Quantity Available to Purchase"`
}
type Result struct {
	Stocks    []*Stock
}

func Worker() string {
	logs.Informational("scanning...")

	url := "https://www.selfridges.com/api/cms/ecom/v1/CN/zh/stock/byId/456-84033258-L8453000"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Informational("err: %+v", err)
		return ""
	}
	req.Header.Set("Api-Key", "xjut2p34999bad9dx7y868ng")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Informational("err: %+v", err)
		return ""
	}

	result := &Result{}
	json.Unmarshal(body, result)

	r := ""
	for _, s := range result.Stocks {
		info := fmt.Sprintf("商品:%s 色号:%s 货号:%s 库存:%s", s.Key, s.Value, s.SkuId, s.Stock)

		if s.Stock != "" {
			st, err := strconv.Atoi(s.Stock)
			if err == nil {
				if st > 0 {
					r = r + info + "\n"
				}
			}
		}
		logs.Informational("info: %+v", info)
	}
	return r
}