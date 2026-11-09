package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type code2SessionResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func code2Session(code string) (*code2SessionResp, error) {
	appID := os.Getenv("WECHAT_APPID")
	secret := os.Getenv("WECHAT_SECRET")
	if appID == "" {
		appID = os.Getenv("APPID")
	}
	if secret == "" {
		secret = os.Getenv("APPSECRET")
	}
	if appID == "" || secret == "" {
		return nil, fmt.Errorf("未配置 WECHAT_APPID / WECHAT_SECRET")
	}

	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	query := url.Values{}
	query.Set("appid", appID)
	query.Set("secret", secret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")

	resp, err := http.Get(endpoint + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result code2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("code2session failed: %d %s", result.ErrCode, result.ErrMsg)
	}
	if result.OpenId == "" {
		return nil, fmt.Errorf("code2session 未返回 openid")
	}
	return &result, nil
}

func openIDFromRequest(r *http.Request) string {
	if id := r.Header.Get("X-WX-OpenId"); id != "" {
		return id
	}
	if id := r.Header.Get("x-wx-openid"); id != "" {
		return id
	}
	return ""
}

func unionIDFromRequest(r *http.Request) string {
	if id := r.Header.Get("X-WX-UnionId"); id != "" {
		return id
	}
	if id := r.Header.Get("x-wx-unionid"); id != "" {
		return id
	}
	return ""
}
