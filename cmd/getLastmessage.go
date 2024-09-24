package cmd

import (
	"fmt"

	"doows/db"
)

// 根据 session_id 检索最新的 last_messages
func GetLastMessages(sessionID int64) (string, error) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return "", fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	var lastMessages string
	query := `
	SELECT last_messages, update_time 
	FROM history_aichat 
	WHERE session_id = ? 
	ORDER BY update_time DESC LIMIT 1`
	err = dbConn.QueryRow(query, sessionID).Scan(&lastMessages)
	if err != nil {
		return "", fmt.Errorf("检索聊天消息时出错: %w", err)
	}

	return lastMessages, nil
}
