package mysql

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	username               string
	password               string
	host                   string
	port                   int
	dbName                 string
	maxOpenConnections     int // 最大打开连接数
	maxIdleConnections     int // 最大空闲连接数
	maxLifetimeConnections int // 连接最大存活时间
	logLevel               int // 日志等级
	slowThreshold          int // 慢查询阈值
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.username, c.password, c.host, c.port, c.dbName)
}

func NewGormDB(ctx context.Context, config Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: config.GetDSN()}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(config.logLevel)),
		})
	if err != nil {
		return nil, fmt.Errorf("connect mysql fail err:%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db fail err:%v", err)
	}
	if config.maxOpenConnections > 0 {
		sqlDB.SetMaxOpenConns(config.maxOpenConnections) //连接池中允许的最大打开（活跃）连接数
	}
	if config.maxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(config.maxIdleConnections) //连接池中允许的最大空闲连接数
	}

	//sqlDB.SetConnMaxIdleTime(time.Minute)                                                // 设置空闲连接的最长时间
	//sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.maxLifetimeConnections)) // 设置连接的最长存活时
	return db, nil
}
