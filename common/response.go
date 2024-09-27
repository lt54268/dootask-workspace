package common

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// 将响应数据序列化为 JSON 并通过 WebSocket 发送给前端
func SendJSONResponse(conn *websocket.Conn, status string, message string) {
	response := Response{
		Status:  status,
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("序列化JSON响应时出错:", err)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"status": "error", "message": "服务器内部错误!"}`))
		return
	}

	conn.WriteMessage(websocket.TextMessage, jsonResponse)
}
