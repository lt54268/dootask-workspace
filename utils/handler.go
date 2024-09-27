package utils

import (
	"doows/cmd"
	"doows/common"
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
	SessionID int64  `json:"sessionId"`
}

type NormalChatRequest struct {
	Slug      string `json:"slug"`
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID int64  `json:"sessionId"`
}

func HandleMessage(conn *websocket.Conn, msg []byte, done chan struct{}) {
	var request WebSocketRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("解析 WebSocket 请求时出错:", err)
		common.SendJSONResponse(conn, "error", "请求格式错误!")
		return
	}

	switch request.Action {
	case "sync":
		log.Println("收到同步用户请求...")
		go func() {
			cmd.SyncUsers(done)
			common.SendJSONResponse(conn, "success", "同步任务已启动!")
		}()

	case "check":
		log.Println("收到检查工作区请求...")
		if err := cmd.CheckWorkspaceLimit(conn); err != nil {
			common.SendJSONResponse(conn, "error", "检查工作区失败!")
		}

	case "get-users":
		log.Println("收到获取已创建工作区用户请求...")
		cmd.GetCreatedWorkspacesUsers(conn)

	case "get-workspace":
		log.Println("收到获取工作区数据请求...")
		var req cmd.GetWorkspaceRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			common.SendJSONResponse(conn, "error", "请求格式错误!")
			return
		}

		wp, err := cmd.GetWorkspacePermission(req.UserID)
		if err != nil {
			common.SendJSONResponse(conn, "error", "获取工作区权限失败:"+err.Error())
			return
		}

		wpJSON, err := json.Marshal(wp)
		if err != nil {
			common.SendJSONResponse(conn, "error", "编码工作区权限数据时出错!")
			return
		}

		conn.WriteMessage(websocket.TextMessage, wpJSON)

	case "set":
		log.Println("收到设置授权请求...")
		var req SetPermissionRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			log.Println("解析权限设置请求时出错:", err)
			common.SendJSONResponse(conn, "error", "权限设置请求格式错误!")
			return
		}

		cmd.SetWorkspacePermission(conn, req.UserID, req.IsCreate)

	case "create":
		log.Println("收到创建工作区请求...")
		var req CreateWorkspaceRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			common.SendJSONResponse(conn, "error", "创建请求格式错误!")
			return
		}

		isAuthorized, err := cmd.CheckUserAuthorization(req.UserID)
		if err != nil {
			common.SendJSONResponse(conn, "error", "检查用户权限时出错!")
			return
		}

		if !isAuthorized {
			common.SendJSONResponse(conn, "error", "用户没有创建工作区的权限!")
			return
		}

		//slug, err := cmd.CreateWorkspace(req.UserID)
		slug, err := cmd.CreateWorkspace(conn, req.UserID)
		if err != nil {
			common.SendJSONResponse(conn, "error", "创建工作区失败: "+err.Error())
			return
		}

		if err := cmd.UpdateWorkspaceID(req.UserID, slug); err != nil {
			common.SendJSONResponse(conn, "error", "更新工作区权限失败: "+err.Error())
			return
		}

		common.SendJSONResponse(conn, "success", "工作区创建成功! slug: "+slug)

	case "stream-chat":
		log.Println("收到流式聊天请求...")
		var req StreamChatRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			sendError(conn, "流式聊天请求格式错误!")
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
			sendError(conn, "常规聊天请求格式错误!")
			return
		}

		res, err := cmd.NormalChat(req.Slug, req.Message, req.Mode, req.SessionID)
		if err != nil {
			sendError(conn, "常规聊天处理失败: "+err.Error())
			return
		}

		sendMessage(conn, res)

	case "insert-message":
		log.Println("收到记录消息请求...")
		var req cmd.BackRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			common.SendJSONResponse(conn, "error", "记录消息请求格式错误!")
			return
		}

		if err := cmd.InsertChatMessage(req); err != nil {
			common.SendJSONResponse(conn, "error", "记录聊天消息失败:"+err.Error())
			return
		}

		common.SendJSONResponse(conn, "success", "消息已成功记录!")

	case "get-history":
		log.Println("收到获取历史记录请求...")
		var req cmd.GetHistoryRequest
		if err := json.Unmarshal(request.Data, &req); err != nil {
			common.SendJSONResponse(conn, "error", "获取请求格式错误!")
			return
		}

		history, err := cmd.GetChatHistory(req.UserID)
		if err != nil {
			common.SendJSONResponse(conn, "error", "获取聊天历史失败: "+err.Error())
			return
		}

		response := cmd.GetHistoryResponse{Messages: history}
		responseData, err := json.Marshal(response)
		if err != nil {
			common.SendJSONResponse(conn, "error", "序列化响应时出错: "+err.Error())
			return
		}

		conn.WriteMessage(websocket.TextMessage, responseData)

	default:
		log.Println("未知请求:", request.Action)
		common.SendJSONResponse(conn, "error", "未知请求类型!")
	}
}
