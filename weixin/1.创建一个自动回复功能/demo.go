package main

import (
	"time"

	"github.com/skynology/cloud-helper"
	"github.com/skynology/cloud-types"
	"github.com/skynology/cloud-types/wechat/mp"
	"github.com/skynology/go-sdk"
)

// 直接回复用户发过来的信息
func Text(h *helper.Helper, req types.CloudRequest, app *skynology.App) {
	// 自动回复
	r, _ := mp.GetText(req.ExtraData)
	res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "您发来了文本信息:\n--------------\n"+r.Content)
	h.Render(res)

}

// 从数据库查询要回复的内容并回复给用户
func TextWithQuery(h *helper.Helper, req types.CloudRequest, app *skynology.App) {
	r, _ := mp.GetText(req.ExtraData)

	// 查询 messages 表有无设置好的自动回复信息
	messages, _, err := app.NewQuery("messages").Equal("text", r.Content).Take(1).Find()
	if err != nil || len(messages) == 0 {
		// 记录log到管理台, 一般用于调试, 默认只保存3天内的log
		if err != nil {
			h.Log(err.String())
		}

		res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "你说了什么?????????")
		h.Render(res)

		// 中途返回时一定记得return语句哦
		return
	}

	content := messages[0].GetString("message")

	res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), content)
	h.Render(res)
}
