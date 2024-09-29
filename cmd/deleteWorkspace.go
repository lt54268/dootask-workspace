package cmd

import (
	"fmt"
	"log"

	"doows/db"
)

// 清空 workspace_id 字段
func DeleteWorkspace(userID int64) error {
	dbConn, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	defer dbConn.Close()

	query := "UPDATE workspace_permission SET workspace_id = NULL WHERE user_id = ?"
	result, err := dbConn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("更新数据库时出错: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("无法获取受影响的行数: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("未找到 user_id 为 %d 的记录", userID)
	}

	log.Printf("用户 %d 的 workspace_id 已成功删除\n", userID)
	return nil
}
