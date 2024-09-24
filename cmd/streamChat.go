package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"doows/db"
)

const streamChatAPI = "http://103.63.139.165:3001/api/v1/workspace/%s/stream-chat"

// RequestPayload 定义发送给API的请求结构
type RequestPayload struct {
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID string `json:"sessionId"`
}

// StreamChat 用于处理流式聊天请求
func StreamChat(slug, message, mode, sessionID string) (string, error) {
	// 构造请求体
	payload := RequestPayload{
		Message:   message,
		Mode:      mode,
		SessionID: sessionID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// 构造请求地址
	url := fmt.Sprintf(streamChatAPI, slug)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 添加请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer VREQHX8-PTGMW06-P9T61XE-BWG31ZW") // 替换为实际的 API Key
	req.Header.Set("accept", "text/event-stream")

	// 发送请求
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var responseData bytes.Buffer
	_, err = io.Copy(&responseData, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return responseData.String(), nil
}

// GetWorkspaceSlug 从 workspace_permission 表中获取 slug
func GetWorkspaceSlug(workspaceID string) (string, error) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return "", fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close() // 确保在函数结束时关闭连接

	var slug string
	query := "SELECT slug FROM workspace_permission WHERE workspace_id = ?"
	err = dbConn.QueryRow(query, workspaceID).Scan(&slug)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve slug: %v", err)
	}
	return slug, nil
}

// // SaveChatSession 将 sessionId 和 slug 存入 history_aichat 表
// func SaveChatSession(sessionID string, slug string, lastMessage string) error {
// 	db, err := db.ConnectDB()
// 	if err != nil {
// 		return fmt.Errorf("数据库连接失败: %w", err)
// 	}
// 	defer db.Close()

// 	query := "INSERT INTO history_aichat (session_id, user_id, last_messages) VALUES (?, ?, ?)"
// 	_, err = db.Exec(query, sessionID, slug, lastMessage) // 使用 slug 作为 user_id
// 	if err != nil {
// 		return fmt.Errorf("保存聊天会话时出错: %w", err)
// 	}

// 	return nil
// }
