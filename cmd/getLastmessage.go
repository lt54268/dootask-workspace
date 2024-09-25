package cmd

import (
	"doows/db"
	"encoding/json"
	"fmt"
)

type ChatMessage struct {
	LastMessage string `json:"last_message"`
}

// 获取指定 session_id 的最新消息并返回 JSON 字符串
func SendChatMessages(sessionID int64) (string, error) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return "", fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	var lastMessage string
	query := `
	SELECT last_messages FROM history_aichat 
	WHERE session_id = ? 
	ORDER BY update_time DESC 
	LIMIT 1`
	err = dbConn.QueryRow(query, sessionID).Scan(&lastMessage)
	if err != nil {
		return "", fmt.Errorf("检索聊天消息时出错: %w", err)
	}

	message := ChatMessage{LastMessage: lastMessage}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("转换为 JSON 时出错: %w", err)
	}

	return string(jsonData), nil
}
