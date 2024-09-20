package db

import (
	"database/sql"
	"doows/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	cfg, err := config.LoadDBConfig("config/config.json")
	if err != nil {
		return nil, fmt.Errorf("无法加载配置文件: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	return db, nil
}
