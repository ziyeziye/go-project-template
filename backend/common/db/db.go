package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var engine *gorm.DB
var mu sync.RWMutex

func initDB() *gorm.DB {
	if engine != nil {
		return engine
	}
	mu.Lock()
	defer mu.Unlock()

	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	dbname := viper.GetString("DB_DATABASE")
	user := viper.GetString("DB_USERNAME")
	pass := viper.GetString("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, pass, host, port, dbname)

	//两个请求同时到来时，都检测到globalDb为nil，其中一个请求会获取锁，然后创建globalDb，并释放锁。另一个请求获取到锁时，globalDb已经可以直接使用了
	if engine != nil {
		return engine
	}

	//记录sql日志
	newLogger := logger.New(
		// log.New(logx.NewMarkLogger("gorm"), "", log.LstdFlags),
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// 禁用 AutoMigrate 自动创建数据库外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic(err)
	// }
	// // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	// sqlDB.SetMaxIdleConns(10)

	// // SetMaxOpenConns 设置打开数据库连接的最大数量。
	// sqlDB.SetMaxOpenConns(100)

	// // SetConnMaxLifetime 设置了连接可复用的最大时间。
	// sqlDB.SetConnMaxLifetime(time.Hour)

	engine = db

	return db
}

func Engine() *gorm.DB {
	if engine != nil {
		return engine
	} else {
		return initDB()
	}
}

func Model(model any) *gorm.DB {
	return Engine().Model(model)
}

func Transaction(fn func(tx *gorm.DB) error) error {
	return Engine().Transaction(fn)
}
