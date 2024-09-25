package cmd

import (
	"doows/db"
	"doows/sync"
	"log"

	"github.com/gorilla/websocket"
)

type UserResponse struct {
	UserID int64 `json:"user_id"`
}

// 获取已经创建工作区的用户数量
func GetCreatedWorkspacesUsers(conn *websocket.Conn) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	getService := sync.NewSyncService(db)

	createdUsers, err := getService.GetCreatedWorkspaces()

	if err != nil {
		log.Println("获取工作区用户失败:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("获取工作区用户失败!"))
		return
	}

	var userResponses []UserResponse
	for _, userID := range createdUsers {
		userResponses = append(userResponses, UserResponse{UserID: userID})
	}

	if err := conn.WriteJSON(userResponses); err != nil {
		log.Println("发送用户列表失败:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("发送用户列表失败!"))
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("获取已创建工作区用户列表完毕!"))
}
