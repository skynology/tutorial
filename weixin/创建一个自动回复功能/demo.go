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

// 用户关注时返回欢迎词
func Subscribe(h *helper.Helper, req types.CloudRequest, app *skynology.App) {
	r, _ := mp.GetSubscribeEvent(req.ExtraData)
	res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "欢迎您关注上空云:)")
	h.Render(res)
}

// 用户关注时保存用户信息并返回欢迎词
func SubscribeWithSave(h *helper.Helper, req types.CloudRequest, app *skynology.App) {
	r, _ := mp.GetSubscribeEvent(req.ExtraData)

	// 若用户从前关注过, 只标记 `subscribe`
	users, _, err := app.NewQuery("_User").Equal("openid", r.FromUserName).Select("nickname", "objectId", "subscribe").Find()
	if err != nil {
		h.Log("query user error:" + err.String())
	}
	if len(users) > 0 {

		// 另一种 go sdk写法
		successed, err := app.NewObjectWithId("_User", users[0].ObjectId).Set("subscribe", 1).Save()
		if !successed {
			h.Log("update user subscribe error:" + err.String())
			return
		}

		// 返回欢迎词
		res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), users[0].GetString("nickname")+", 欢迎您再次关注上空云:)")
		h.Render(res)
		return
	}

	// 从微信服务器获取用户信息
	user, err := app.GetWeixin("users/" + r.FromUserName)
	if err != nil {
		h.Log("get user error:" + err.String())
		return
	}

	// 保存微信用户信息到云平台
	// 因用户表的username 及 password为必填字段, 需设置值
	user["username"] = r.FromUserName
	user["password"] = "123456"

	successed, err := app.NewObject("_User").SetMulti(user).Save()
	if !successed {
		if err != nil {
			h.Log("create user error:" + err.String())
		}
		return
	}

	// 返回欢迎词
	res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "欢迎您关注上空云:)")
	h.Render(res)

}

// 用户取消关注时标记 `subscribe`
func UnSubscribeWithSave(h *helper.Helper, req types.CloudRequest, app *skynology.App) {
	r, _ := mp.GetUnsubscribeEvent(req.ExtraData)

	users, _, err := app.NewQuery("_User").Equal("openid", r.FromUserName).Find()
	if err != nil {
		h.Log("get user error:" + err.String())
		return
	}

	if len(users) == 0 {
		return
	}

	user := users[0]

	// 把关注了又取消的用户的 `subscribe` 标记为 2
	successed, err := user.Set("subscribe", 2).Save()
	if !successed && err != nil {
		h.Log(err.String())
	}

}
