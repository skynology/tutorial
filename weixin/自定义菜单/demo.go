package main

import (
	"fmt"
	"time"

	"github.com/skynology/cloud-helper"
	"github.com/skynology/cloud-types"
	"github.com/skynology/cloud-types/wechat/mp"
	"github.com/skynology/go-sdk"
)

// 点击自定义菜单按钮: (Click类型)
func MenuClick(h *helper.Helper, req types.CloudRequest, app *skynology.App) {

	r, _ := mp.GetClickEvent(req.ExtraData)

	if r.EventKey == "contactus" {
		res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "请加QQ群: 434581608")
		h.Render(res)
		return
	}

	// 点了赞的话, 更新系统 'test' 表中的被赞的数量
	if r.EventKey == "zan" {

		// 这里直接更新了, 未做判断是否成功
		app.NewObjectWithId("test", "55549e983a859a7649000001").Increment("zan").Save()

		// 重新查询下, 并返回当前赞的数
		result, _ := app.NewQuery("test").GetObject("55549e983a859a7649000001")
		content := fmt.Sprintf("当前被赞了%v次!", result.Get("zan"))
		res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), content)
		h.Render(res)
	}

	if r.EventKey == "img" {
		images, _, _ := app.NewQuery("image").Take(1).Find()

		// 可按发起人来查询
		//images, _, _ := app.NewQuery("image").Equal("openid", r.FromUserName).Take(1).Find()

		// 当用户从未发送图片时, 直接返回提示语
		if len(images) == 0 {
			res := mp.NewResText(r.FromUserName, r.ToUserName, time.Now().Unix(), "您还从未发送过图片哦^_^")
			h.Render(res)
			return
		}

		res := mp.NewResImage(r.FromUserName, r.ToUserName, time.Now().Unix(), images[0].GetString("mediaId"))
		h.Render(res)
		return
	}
}

// 用户点击拍照后发送图片过来
func SendImage(h *helper.Helper, req types.CloudRequest, app *skynology.App) {

	r, _ := mp.GetImage(req.ExtraData)

	// 直接返回用户发过来的图片
	res := mp.NewResImage(r.FromUserName, r.ToUserName, time.Now().Unix(), r.MediaId)
	h.Render(res)

	// meidaId及用户的openid 保存到服务器
	app.NewObject("image").Set("mediaId", r.MediaId).Set("openid", r.FromUserName).Save()

}
