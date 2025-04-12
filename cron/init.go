package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

// Init 初始化定时任务
func Init() error {
	log.Println("开始进行初始化和定时任务")

	c = cron.New(cron.WithSeconds())

	// 每分钟检查一次需要同步的API
	_, err := c.AddFunc("0 * * * * *", func() {
		if err := checkAndSyncAPIs(); err != nil {
			log.Printf("检查并同步API数据失败: %v", err)
		}
	})
	if err != nil {
		return err
	}
	// 启动定时任务
	c.Start()
	log.Println("定时任务初始化成功")
	return nil
}

// checkAndSyncAPIs 检查并同步需要更新的API数据
func checkAndSyncAPIs() error {
	return nil
}

// Stop 停止定时任务
func Stop() {
	if c != nil {
		c.Stop()
	}
}
