package sync

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type SyncService struct {
	db *sql.DB
}

func NewSyncService(db *sql.DB) *SyncService {
	return &SyncService{db: db}
}

// 同步 pre_users 表和 workspace_permission 表
func (s *SyncService) SyncUsers() error {
	preUserIDs := make(map[int64]bool)
	rows, err := s.db.Query("SELECT userid FROM pre_users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		preUserIDs[userID] = true

		var exists bool
		err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM workspace_permission WHERE user_id = ?)", userID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err := s.db.Exec("INSERT INTO workspace_permission (user_id, is_create, create_time, update_time) VALUES (?, ?, NOW(), NOW())", userID, false)
			if err != nil {
				return err
			}
			log.Printf("Inserted user_id: %d into workspace_permission\n", userID)
		}

		// workspace_permission 表同步 user_id 到 history_aichat 表
		// err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM history_aichat WHERE user_id = ?)", userID).Scan(&exists)
		// if err != nil {
		// 	return err
		// }

		// if !exists {
		// 	_, err := s.db.Exec("INSERT INTO history_aichat (session_id, model, user_id, last_messages, create_time, update_time) VALUES (?, ?, ?, ?, NOW(), NOW())", 0, 0, userID, "")
		// 	if err != nil {
		// 		return err
		// 	}
		// 	log.Printf("Inserted user_id: %d into history_aichat\n", userID)
		// }
	}

	if err := rows.Err(); err != nil {
		return err
	}

	placeholders := generatePlaceholders(len(preUserIDs))
	args := make([]interface{}, len(preUserIDs))
	i := 0
	for id := range preUserIDs {
		args[i] = id
		i++
	}

	_, err = s.db.Exec("DELETE FROM workspace_permission WHERE user_id NOT IN ("+placeholders+")", args...)
	if err != nil {
		return err
	}

	//log.Println("pre_users 表用户已同步至 workspace_permission 表!")
	return nil
}

func generatePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

// 检查创建工作区数量是否超过 3
func (s *SyncService) CheckWorkspaceCreationLimit() error {
	//log.Println("开始检查已创建工作区的数量...")

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM workspace_permission WHERE workspace_id IS NOT NULL").Scan(&count)
	if err != nil {
		return err
	}

	// 检查数量是否超过 3
	if count >= 3 {
		log.Printf("Warning: 超过工作区创建的数量，当前数量: %d", count)
		return errors.New("工作区创建数量已达最大!")
	}

	log.Printf("当前已创建工作区数量: %d", count)
	return nil
}

// 查询 workspace_id 不为空的所有用户
func (s *SyncService) GetCreatedWorkspaces() ([]int64, error) {
	var userIds []int64

	rows, err := s.db.Query("SELECT user_id FROM workspace_permission WHERE workspace_id IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIds = append(userIds, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}

// 检查并设置 workspace_permission 表中的 is_create 状态和 workspace_id
func (s *SyncService) SetWorkspacePermission(userID int64, isCreate bool) error {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pre_users WHERE userid = ?)`
	err := s.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("数据库查询失败: %v", err)
	}

	if !exists {
		return fmt.Errorf("该用户不存在: %d", userID)
	}

	query = `SELECT EXISTS(SELECT 1 FROM workspace_permission WHERE user_id = ?)`
	err = s.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("数据库查询失败: %v", err)
	}

	now := time.Now()

	if exists {
		// 如果存在，则更新 is_create 和 update_time
		updateQuery := `UPDATE workspace_permission SET is_create = ?, update_time = ? WHERE user_id = ?`
		_, err = s.db.Exec(updateQuery, isCreate, now, userID)
		if err != nil {
			return fmt.Errorf("更新工作区权限失败: %v", err)
		}
		log.Printf("用户 %d 的工作区权限更新成功!", userID)
	} else {
		// 如果不存在，则插入一条新记录
		insertQuery := `INSERT INTO workspace_permission (user_id, is_create, create_time, update_time) VALUES (?, ?, ?, ?)`
		_, err = s.db.Exec(insertQuery, userID, isCreate, now, now)
		if err != nil {
			return fmt.Errorf("插入工作区权限失败: %v", err)
		}
		log.Printf("用户 %d 的工作区权限插入成功!", userID)
	}

	return nil
}
