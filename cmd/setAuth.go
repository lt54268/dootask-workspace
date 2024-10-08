package cmd

import (
	"doows/common"
	"doows/db"
	"doows/sync"
	"log"

	"github.com/gorilla/websocket"
)

// 设置指定用户的 is_create 状态和 workspace_id
func SetWorkspacePermission(conn *websocket.Conn, userID int64, isCreate bool) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
		common.SendJSONResponse(conn, "error", "数据库连接失败!")
		return
	}
	defer db.Close()

	setService := sync.NewSyncService(db)

	// 只有在设置为 true 时检查当前的工作区数量是否超过限定值
	if isCreate {
		err = setService.CheckWorkspaceCreationLimit()
		if err != nil {
			log.Println("工作区创建检查时出错:", err)
			common.SendJSONResponse(conn, "error", "工作区创建检查时出错!")
			return
		}

		if err != nil && err.Error() == "工作区创建数量已达最大!" {
			common.SendJSONResponse(conn, "error", "创建工作区数量已达最大!")
			return
		}
	}

	if err := setService.SetWorkspacePermission(userID, isCreate); err != nil {
		log.Println("设置工作区权限失败:", err)
		common.SendJSONResponse(conn, "error", "设置工作区权限失败!"+err.Error())
		return
	}

	common.SendJSONResponse(conn, "success", "工作区权限设置成功!")
}
