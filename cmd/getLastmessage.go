package cmd

import (
	"doows/db"
	"fmt"
	"time"
)

type GetHistoryRequest struct {
	UserID int64 `json:"user_id"`
}

type ChatMessage struct {
	ID           int64     `json:"id"`
	SessionID    int64     `json:"session_id"`
	Model        string    `json:"model"`
	UserID       int64     `json:"user_id"`
	LastMessages string    `json:"last_messages"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type GetHistoryResponse struct {
	Messages []ChatMessage `json:"messages"`
}

// 获取指定 session_id 的最新消息并返回 JSON 字符串
func GetChatHistory(userID int64) ([]ChatMessage, error) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	query := `SELECT id, session_id, model, user_id, last_messages, create_time, update_time
	          FROM history_aichat
	          WHERE user_id = ?`

	rows, err := dbConn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询数据库时出错: %w", err)
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		var createTimeStr string
		var updateTimeStr string
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Model, &msg.UserID, &msg.LastMessages, &createTimeStr, &updateTimeStr); err != nil {
			return nil, fmt.Errorf("解析查询结果时出错: %w", err)
		}

		if err := parseAndSetTime(createTimeStr, &msg.CreateTime); err != nil {
			return nil, err
		}
		if err := parseAndSetTime(updateTimeStr, &msg.UpdateTime); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("处理查询结果时出错: %w", err)
	}

	return messages, nil
}

func parseAndSetTime(timeStr string, target *time.Time) error {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("解析时间时出错: %w", err)
	}
	*target = parsedTime
	return nil
}
