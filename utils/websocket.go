package utils

import (
	"doows/cmd"
	"log"
	"net/http"
	stdsync "sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(conn *websocket.Conn, done chan struct{}) {
	defer conn.Close()

	var wg stdsync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		cmd.SyncUsers(done)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}

		HandleMessage(conn, msg, done)
	}

	wg.Wait()
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, done chan struct{}) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级连接失败:", err)
		return
	}

	handleConnections(conn, done)
}

// sendMessage 发送成功消息到 WebSocket
func sendMessage(conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("发送 WebSocket 消息失败:", err)
	}
}

// sendError 发送错误消息到 WebSocket
func sendError(conn *websocket.Conn, errorMsg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte("错误: "+errorMsg)); err != nil {
		log.Println("发送 WebSocket 错误消息失败:", err)
	}
}
