package cmd

import (
	"bytes"
	"doows/db"
	"encoding/json"
	"fmt"
	"net/http"
)

// 创建工作区并返回 slug
func CreateWorkspace(userID int64) (string, error) {
	// 创建 HTTP 请求
	url := "http://103.63.139.165:3001/api/v1/workspace/new"
	body := map[string]string{"name": fmt.Sprintf("Workspace for User %d", userID)}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("json.Marshal 错误: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求时出错: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer VREQHX8-PTGMW06-P9T61XE-BWG31ZW") // 替换为实际的 API Key
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求时出错: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var response struct {
		Workspace struct {
			Slug string `json:"slug"`
		} `json:"workspace"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应时出错: %w", err)
	}

	return response.Workspace.Slug, nil
}

// 更新 workspace_permission 表中的 slug
func UpdateWorkspaceID(userID int64, slug string) error {
	db, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	// 这里是更新数据库的逻辑，假设你有对应的数据库操作
	// 示例：使用 SQL 语句更新 workspace_permission 表
	query := "UPDATE workspace_permission SET workspace_id = ? WHERE user_id = ?"
	_, err = db.Exec(query, slug, userID)
	if err != nil {
		return fmt.Errorf("更新数据库时出错: %w", err)
	}

	return nil
}
