package router

import (
	"doows/config"
	"doows/utils"
	"log"
	"net/http"
)

func RunWebSocket() {
	cfg := config.LoadWBConfig()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		utils.HandleWebSocket(w, r, done)
	})

	log.Printf("WebSocket 服务器启动,监听端口:%s", cfg.WebSocketPort)

	if err := http.ListenAndServe(":"+cfg.WebSocketPort, nil); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
