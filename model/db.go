package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"omo-msa-license/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/micro/go-micro/v2/logger"
	uuid "github.com/satori/go.uuid"
)

var base64Coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

var dialectSqlName string
var dialectSqlArgs string

func Setup() {
	if config.Schema.Database.Lite {
		dialectSqlName = "sqlite3"
		sqlite_filepath := config.Schema.Database.SQLite.Path
		dialectSqlArgs = sqlite_filepath
		logger.Warnf("!!! Database is lite mode, file at %v", sqlite_filepath)
	} else {
		dialectSqlName = "mysql"
		mysql_addr := config.Schema.Database.MySQL.Address
		mysql_user := config.Schema.Database.MySQL.User
		mysql_passwd := config.Schema.Database.MySQL.Password
		mysql_db := config.Schema.Database.MySQL.DB
		dialectSqlArgs = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True", mysql_user, mysql_passwd, mysql_addr, mysql_db)
	}
}

func AutoMigrateDatabase() {
	db, err := openSqlDB()
	if nil != err {
		panic(err)
	}
	defer closeSqlDB(db)

	err = db.AutoMigrate(&Space{}).Error
	if nil != err {
		panic(err)
	}
	err = db.AutoMigrate(&Key{}).Error
	if nil != err {
		panic(err)
	}
	err = db.AutoMigrate(&Certificate{}).Error
	if nil != err {
		panic(err)
	}
}

func openSqlDB() (*gorm.DB, error) {
	return gorm.Open(dialectSqlName, dialectSqlArgs)
}

func closeSqlDB(_db *gorm.DB) {
	_db.Close()
}

func NewUUID() string {
	guid := uuid.NewV4()
	h := md5.New()
	h.Write(guid.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func ToUUID(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(_content string) string {
	h := md5.New()
	h.Write([]byte(_content))
	return hex.EncodeToString(h.Sum(nil))
}

func ToBase64(_content []byte) string {
	return base64Coder.EncodeToString(_content)
}
