package cmd

import (
	"doows/db"
	"doows/sync"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 定时同步用户信息
func SyncUsers(done chan struct{}) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	userService := sync.NewSyncService(db)

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//log.Println("开始同步用户...")
			if err := userService.SyncUsers(); err != nil {
				log.Println("同步用户失败:", err)
			} else {
				log.Println("同步完成!")
			}
		case <-done:
			log.Println("停止用户同步任务...")
			return
		}
	}
}
