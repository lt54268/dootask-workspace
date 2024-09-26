package cmd

import (
	"fmt"
	"time"

	"doows/db"
)

type BackRequest struct {
	SessionID   int64  `json:"session_id"`
	UserID      int64  `json:"user_id"`
	LastMessage string `json:"last_message"`
	model       string `json:"model"`
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
	INSERT INTO history_aichat (session_id, model, user_id, last_messages, create_time, update_time) 
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE last_messages = ?, update_time = ?`
	_, err = dbConn.Exec(query, req.SessionID, req.model, req.UserID, req.LastMessage, currentTime, currentTime, req.LastMessage, currentTime)
	if err != nil {
		return fmt.Errorf("保存聊天消息时出错: %w", err)
	}

	return nil
}
