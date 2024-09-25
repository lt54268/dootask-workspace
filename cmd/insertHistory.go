package cmd

import (
	"fmt"
	"time"

	"doows/db"
)

type BackRequest struct {
	SessionID   int64  `json:"session_id"`
	Slug        string `json:"slug"` // 存放 user_id
	LastMessage string `json:"last_message"`
}

// 将前端发送的对话数据插入到 history_aichat 表中
func InsertChatMessage(req BackRequest) error {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	currentTime := time.Now()

	query := `
	INSERT INTO history_aichat (session_id, user_id, last_messages, create_time, update_time) 
	VALUES (?, ?, ?, ?, ?)`
	_, err = dbConn.Exec(query, req.SessionID, req.Slug, req.LastMessage, currentTime, currentTime)
	if err != nil {
		return fmt.Errorf("保存聊天消息时出错: %w", err)
	}

	return nil
}
