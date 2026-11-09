package service

import (
	"encoding/json"
	"net/http"
)

// JsonResult 统一返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

func writeJSON(w http.ResponseWriter, res *JsonResult) {
	msg, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func writeOK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, &JsonResult{Code: 0, Data: data})
}

func writeErr(w http.ResponseWriter, msg string) {
	writeJSON(w, &JsonResult{Code: -1, ErrorMsg: msg})
}

// HealthHandler 健康检查，用于确认新版本已部署
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]interface{}{
		"ok":      true,
		"service": "taroGame",
		"apis":    []string{"/api/health", "/api/login", "/api/user", "/api/count"},
	})
}
