package router

import (
	"doows/config"
	"log"
	"net/http"
)

func RunHttp() {
	// 加载 HTTP 端口配置
	cfg := config.LoadWBConfig()
	port := cfg.WebSocketPort

	// // 注册 HTTP 路由
	// http.HandleFunc("/sync", )
	// http.HandleFunc("/check", )
	// http.HandleFunc("/get", )

	// 启动 HTTP 服务器
	log.Printf("HTTP 服务器启动中,监听端口:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("HTTP 服务器启动失败:", err)
	}
}
