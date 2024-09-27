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

const NormalChatAPI = "http://103.63.139.165:3001/api/v1/workspace/%s/chat"

type RequestPayloadNormal struct {
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID int64  `json:"sessionId"`
}

// 处理常规聊天请求
func NormalChat(slug, message string, mode string, sessionID int64) (string, error) {
	payload := RequestPayloadNormal{
		Message:   message,
		Mode:      mode,
		SessionID: sessionID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("调用有效载荷失败: %v", err)
	}

	url := fmt.Sprintf(NormalChatAPI, slug)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer VREQHX8-PTGMW06-P9T61XE-BWG31ZW") // 替换为实际的 API Key
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	var responseData bytes.Buffer
	_, err = io.Copy(&responseData, resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return responseData.String(), nil
}

// 从 workspace_permission 表中获取 slug
func GetWorkspaceSlugnomal(workspaceID string) (string, error) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return "", fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	var slug string
	query := "SELECT slug FROM workspace_permission WHERE workspace_id = ?"
	err = dbConn.QueryRow(query, workspaceID).Scan(&slug)
	if err != nil {
		return "", fmt.Errorf("未能检索到slug: %v", err)
	}
	return slug, nil
}
