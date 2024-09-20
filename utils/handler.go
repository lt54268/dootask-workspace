package utils

import (
	"doows/cmd"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketRequest struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type SetPermissionRequest struct {
	UserID      int64  `json:"user_id"`
	IsCreate    bool   `json:"is_create"`
	WorkspaceID string `json:"workspace_id"`
}

// 处理 WebSocket 消息
func HandleMessage(conn *websocket.Conn, msg []byte, done chan struct{}) {
	var request WebSocketRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("解析 WebSocket 请求时出错:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("请求格式错误"))
		return
	}

	switch request.Action {
	case "sync":
		log.Println("收到同步请求...")
		go func() {
			cmd.SyncUsers(done)
			conn.WriteMessage(websocket.TextMessage, []byte("同步任务已启动!"))
		}()

	case "check":
		log.Println("收到检查请求...")
		cmd.CheckWorkspaceLimit(conn)

	case "get":
		log.Println("收到获取请求...")
		cmd.GetCreatedWorkspacesUsers(conn)

	case "set":
		log.Println("收到设置权限请求...")
		var req SetPermissionRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			log.Println("解析权限设置请求时出错:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("权限设置请求格式错误"))
			return
		}
		if err := cmd.CheckWorkspaceLimit(conn); err != nil {
			log.Println("工作区创建限制检查失败:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("工作区创建数量已达最大"))
			return
		}
		cmd.SetWorkspacePermission(conn, req.UserID, req.IsCreate, req.WorkspaceID)

	default:
		log.Println("未知请求:", request.Action)
		conn.WriteMessage(websocket.TextMessage, []byte("未知请求类型"))
	}
}
