package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // 纯Go实现的SQLite驱动
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	// 获取数据库目录（优先使用环境变量中指定的数据目录）
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	log.Printf("数据目录: %s", dataDir)

	// 确保data目录存在
	if err := os.MkdirAll(dataDir, 0777); err != nil {
		log.Printf("创建数据目录失败: %v", err)
		return err
	}

	// 数据库文件路径
	dbPath := filepath.Join(dataDir, "database.db")
	log.Printf("使用数据库文件: %s", dbPath)

	// 如果数据库文件不存在，并且有写权限，则尝试创建一个空文件
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		emptyFile, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Printf("创建数据库文件失败: %v", err)
			return err // 如果创建失败，直接返回错误而不继续
		} else {
			emptyFile.Close()
			// 确保文件权限设置正确
			if err := os.Chmod(dbPath, 0666); err != nil {
				log.Printf("设置数据库文件权限失败: %v", err)
				// 继续执行，因为文件已创建
			}
			log.Printf("创建了空数据库文件")
		}
	}

	// 检查文件是否可写
	if file, err := os.OpenFile(dbPath, os.O_WRONLY, 0666); err != nil {
		log.Printf("警告：数据库文件不可写: %v", err)
		// 尝试修复权限
		if err := os.Chmod(dbPath, 0666); err != nil {
			log.Printf("无法修改数据库文件权限: %v", err)
		}
	} else {
		file.Close()
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Printf("连接数据库失败: %v", err)
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}
