package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type WBConfig struct {
	WebSocketPort string
}

func LoadDBConfig(filePath string) (*DBConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config DBConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadWBConfig() *WBConfig {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3333" // 默认端口
	}

	return &WBConfig{
		WebSocketPort: port,
	}
}
