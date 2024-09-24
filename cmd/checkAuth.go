package cmd

import (
	"database/sql"
	"doows/db"
)

// 检查用户是否有创建工作区的权限
func CheckUserAuthorization(userID int64) (bool, error) {
	dbConn, err := db.ConnectDB() // 假设 db.ConnectDB() 返回数据库连接
	if err != nil {
		return false, err
	}
	defer dbConn.Close()

	var isCreate bool
	query := "SELECT is_create FROM workspace_permission WHERE user_id = ?"
	err = dbConn.QueryRow(query, userID).Scan(&isCreate)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果用户不在表中，返回 false
			return false, nil
		}
		return false, err
	}

	return isCreate, nil
}
