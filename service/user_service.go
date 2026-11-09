package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type loginRequest struct {
	Code      string `json:"code"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
}

// LoginHandler 微信登录：写入/更新用户
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, "仅支持 POST")
		return
	}

	var req loginRequest
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&req)
		defer r.Body.Close()
	}

	openId := openIDFromRequest(r)
	unionId := unionIDFromRequest(r)

	if openId == "" && strings.TrimSpace(req.Code) != "" {
		session, err := code2Session(strings.TrimSpace(req.Code))
		if err != nil {
			writeErr(w, err.Error())
			return
		}
		openId = session.OpenId
		if session.UnionId != "" {
			unionId = session.UnionId
		}
	}

	if openId == "" {
		writeErr(w, "未获取到微信 openId，请通过云托管调用或传入 code")
		return
	}

	user := &model.UserModel{
		OpenId:    openId,
		UnionId:   unionId,
		NickName:  strings.TrimSpace(req.NickName),
		AvatarUrl: strings.TrimSpace(req.AvatarUrl),
	}

	saved, err := dao.UserImp.UpsertByOpenId(user)
	if err != nil {
		writeErr(w, err.Error())
		return
	}
	writeOK(w, saved)
}

// ProfileHandler 获取当前用户信息
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, "仅支持 GET")
		return
	}

	openId := openIDFromRequest(r)
	if openId == "" {
		openId = strings.TrimSpace(r.URL.Query().Get("openId"))
	}
	if openId == "" {
		writeErr(w, "未获取到微信 openId")
		return
	}

	user, err := dao.UserImp.GetByOpenId(openId)
	if err != nil {
		if dao.IsNotFound(err) {
			writeErr(w, "用户不存在，请先登录")
			return
		}
		writeErr(w, err.Error())
		return
	}
	writeOK(w, user)
}
