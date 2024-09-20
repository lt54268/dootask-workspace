package cmd

import (
	"doows/db"
	"doows/sync"
	"log"

	"github.com/gorilla/websocket"
)

// 设置指定用户的 is_create 状态和 workspace_id
func SetWorkspacePermission(conn *websocket.Conn, userID int64, isCreate bool, workspaceID string) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("数据库连接失败!"))
		return
	}
	defer db.Close()

	setService := sync.NewSyncService(db)

	// 检查当前 is_create = true 的用户数量是否超过限定值
	err = setService.CheckWorkspaceCreationLimit()
	if err != nil {
		log.Println("工作区创建检查时出错:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("工作区创建检查时出错!"))
		return
	}

	// 在不为 nil 时再检查错误消息
	if err != nil && err.Error() == "工作区创建数量已达最大!" {
		conn.WriteMessage(websocket.TextMessage, []byte("创建工作区数量已达最大!"))
		return
	}

	// 设置工作区权限
	if err := setService.SetWorkspacePermission(userID, isCreate, workspaceID); err != nil {
		log.Println("设置工作区权限失败:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("设置工作区权限失败: "+err.Error()))
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("工作区权限设置成功!"))
}