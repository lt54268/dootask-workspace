package cmd

import (
	"database/sql"
	"doows/common"
	"doows/db"
	"doows/sync"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type UserResponse struct {
	UserID int64 `json:"user_id"`
}

type GetWorkspaceRequest struct {
	UserID int64 `json:"user_id"`
}

type WorkspacePermission struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	IsCreate    bool      `json:"is_create"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
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
		common.SendJSONResponse(conn, "error", "获取工作区用户失败!")
		return
	}

	var userResponses []UserResponse
	for _, userID := range createdUsers {
		userResponses = append(userResponses, UserResponse{UserID: userID})
	}

	if err := conn.WriteJSON(userResponses); err != nil {
		log.Println("发送用户列表失败:", err)
		common.SendJSONResponse(conn, "error", "发送用户列表失败!")
		return
	}

	//common.SendJSONResponse(conn, "success", "获取已创建工作区用户列表完毕!")
}

// 根据 user_id 获取 workspace_permission 表中的所有字段
func GetWorkspacePermission(userID int64) (*WorkspacePermission, error) {
	db, err := db.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	var wp WorkspacePermission
	var workspaceID sql.NullString
	var createTime, updateTime []byte
	query := `SELECT id, user_id, is_create, workspace_id, create_time, update_time 
              FROM workspace_permission WHERE user_id = ?`
	err = db.QueryRow(query, userID).Scan(
		&wp.ID, &wp.UserID, &wp.IsCreate, &workspaceID, &createTime, &updateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("未找到用户 %d 的工作区权限记录", userID)
		}
		return nil, fmt.Errorf("查询数据库时出错: %w", err)
	}

	wp.CreateTime, err = parseTime(createTime)
	if err != nil {
		return nil, fmt.Errorf("解析创建时间时出错: %w", err)
	}
	wp.UpdateTime, err = parseTime(updateTime)
	if err != nil {
		return nil, fmt.Errorf("解析更新时间时出错: %w", err)
	}

	if workspaceID.Valid {
		wp.WorkspaceID = workspaceID.String
	}

	return &wp, nil
}

func parseTime(timeBytes []byte) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, string(timeBytes))
}
