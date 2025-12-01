package dao

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var _db *gorm.DB

func Database(connRead, connWrite string) {
	var ormlogger logger.Interface
	if gin.Mode() == "debug" {
		ormlogger = logger.Default.LogMode(logger.Info)
	} else {
		ormlogger = logger.Default
	}
	// 增加重试机制，避免 MySQL 尚未完全启动时直接 panic
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       connRead,
			DefaultStringSize:         256,  // String类型默认字段长度
			DisableDatetimePrecision:  true, // 禁用 datetime 精度，mysql 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true, // 重新命名索引，就要把索引先删除再重建，mysql 5.7 不支持
			DontSupportRenameColumn:   true, // 用 change 重命名列，mysql 8 之前的数据库不支持
			SkipInitializeWithVersion: false,
		}), &gorm.Config{
			Logger: ormlogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err == nil {
			break
		}
		// 等待 MySQL 服务就绪
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)  // 设置连接池
	sqlDB.SetMaxOpenConns(100) // 打开连接数
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	_db = db

	// 主从配置
	_ = _db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(connWrite)}, // 写操作
		Replicas: []gorm.Dialector{mysql.Open(connRead)},  // 读操作
		Policy:   dbresolver.RandomPolicy{},
	}))

	Migration()
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
