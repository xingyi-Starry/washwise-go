package model

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接并执行自动迁移
// path: SQLite 数据库文件路径
func InitDB(path string) error {
	// 确保数据库文件所在目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 打开数据库连接
	var err error
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&Machine{}); err != nil {
		return err
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
