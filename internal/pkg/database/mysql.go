package database

import (
	"database/sql"
	"fmt"
	"go-gen/config"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	// 引用数据库驱动初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	maxOpenConns    = 150
	maxIdleConns    = 100
	connMaxLifetime = 100
)

var engine *gorm.DB
var dbName string = "damn"

// GetDB get gorm.DB
func GetDB() *gorm.DB {
	return engine
}

// GetDbName get database name
func GetDbName() string {
	return dbName
}

// Close closes current db connection
func Close() error {
	return engine.Close()
}

// SetUp set up db connection
func SetUp() {
	var err error
	mysqlConf := config.EnvConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		mysqlConf.User,
		mysqlConf.Pwd,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.DBName,
		mysqlConf.Charset,
		true,
		"Local")

	engine, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("connect to mysql fail, ", dsn, err)
		panic(err)
	}
	engine.LogMode(true)

	engine.DB().SetConnMaxLifetime(connMaxLifetime * time.Second)
	engine.DB().SetMaxOpenConns(maxOpenConns)
	engine.DB().SetMaxIdleConns(maxIdleConns)

	// 禁止update/delete传空对象
	engine.BlockGlobalUpdate(true)
}

// MockDB
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	engine, err = gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}
	engine.LogMode(true)
	return db, mock
}
