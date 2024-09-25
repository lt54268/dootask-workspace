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
	UserID   int64 `json:"user_id"`
	IsCreate bool  `json:"is_create"`
}

type CreateWorkspaceRequest struct {
	UserID int64 `json:"user_id"`
}

type StreamChatRequest struct {
	Slug      string `json:"slug"`
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID string `json:"sessionId"`
}

type NormalChatRequest struct {
	Slug      string `json:"slug"`
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID string `json:"sessionId"`
}

func HandleMessage(conn *websocket.Conn, msg []byte, done chan struct{}) {
	var request WebSocketRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("解析 WebSocket 请求时出错:", err)
		sendError(conn, "请求格式错误")
		return
	}

	switch request.Action {
	case "sync":
		log.Println("收到同步请求...")
		go func() {
			cmd.SyncUsers(done)
			sendMessage(conn, "同步任务已启动!")
		}()

	case "check":
		log.Println("收到检查请求...")
		if err := cmd.CheckWorkspaceLimit(conn); err != nil {
			sendError(conn, "检查失败: "+err.Error())
		}

	case "get":
		log.Println("收到获取请求...")
		cmd.GetCreatedWorkspacesUsers(conn)

	case "set":
		log.Println("收到设置权限请求...")
		var req SetPermissionRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			log.Println("解析权限设置请求时出错:", err)
			sendError(conn, "权限设置请求格式错误")
			return
		}
		// if err := cmd.CheckWorkspaceLimit(conn); err != nil {
		// 	sendError(conn, "工作区创建数量已达最大!")
		// 	return
		// }
		cmd.SetWorkspacePermission(conn, req.UserID, req.IsCreate)

	case "create":
		log.Println("收到创建工作区请求...")
		var req CreateWorkspaceRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			sendError(conn, "创建请求格式错误")
			return
		}

		isAuthorized, err := cmd.CheckUserAuthorization(req.UserID)
		if err != nil {
			sendError(conn, "检查用户权限时出错")
			return
		}

		if !isAuthorized {
			sendError(conn, "用户没有创建工作区的权限!")
			return
		}

		slug, err := cmd.CreateWorkspace(req.UserID)
		if err != nil {
			sendError(conn, "创建工作区失败: "+err.Error())
			return
		}

		if err := cmd.UpdateWorkspaceID(req.UserID, slug); err != nil {
			sendError(conn, "更新工作区权限失败: "+err.Error())
			return
		}

		sendMessage(conn, "工作区创建成功! slug: "+slug)

	case "stream-chat":
		log.Println("收到流式聊天请求...")
		var req StreamChatRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			sendError(conn, "流式聊天请求格式错误")
			return
		}

		response, err := cmd.StreamChat(req.Slug, req.Message, req.Mode, req.SessionID)
		if err != nil {
			sendError(conn, "流式聊天处理失败: "+err.Error())
			return
		}

		sendMessage(conn, response)

	case "chat":
		log.Println("收到常规聊天请求...")
		var req NormalChatRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			sendError(conn, "常规聊天请求格式错误")
			return
		}

		res, err := cmd.NormalChat(req.Slug, req.Message, req.Mode, req.SessionID)
		if err != nil {
			sendError(conn, "常规聊天处理失败: "+err.Error())
			return
		}

		sendMessage(conn, res)

	case "back":
		log.Println("收到返回消息请求...")
		var req cmd.BackRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			sendError(conn, "返回消息请求格式错误")
			return
		}

		if err := cmd.InsertChatMessage(req); err != nil {
			sendError(conn, "插入聊天消息失败: "+err.Error())
			return
		}

		sendMessage(conn, "消息已成功记录!")

	case "send":
		log.Println("收到发送消息请求...")
		var sessionID int64
		if err := json.Unmarshal(request.Data, &sessionID); err != nil {
			sendError(conn, "发送消息请求格式错误")
			return
		}

		lastMessage, err := cmd.SendChatMessages(sessionID)
		if err != nil {
			sendError(conn, "获取聊天消息失败: "+err.Error())
			return
		}

		sendMessage(conn, lastMessage)

	default:
		log.Println("未知请求:", request.Action)
		sendError(conn, "未知请求类型")
	}
}
