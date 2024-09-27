package cmd

import (
	"doows/common"
	"doows/db"
	"doows/sync"
	"log"

	"github.com/gorilla/websocket"
)

// 检查获授权的用户数量是否超过限制
func CheckWorkspaceLimit(conn *websocket.Conn) error {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	checkService := sync.NewSyncService(db)

	if err := checkService.CheckWorkspaceCreationLimit(); err != nil {
		log.Println("工作区创建检查时出错:", err)
		return err
	}

	log.Println("工作区创建限制检查完毕!")
	common.SendJSONResponse(conn, "success", "工作区创建限制检查完毕!")
	return nil
}
