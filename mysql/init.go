package mysql

import (
	"fmt"
	"generate/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var mysqlOnce sync.Once
var db *gorm.DB
var err error

func InitMysql() *gorm.DB {
	if db != nil {
		return db
	}

	mysqlOnce.Do(func() {
		baseConf := config.YamlConfig.MysqlConfig
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			baseConf.User,
			baseConf.Password,
			baseConf.Host,
			baseConf.Port,
			baseConf.DBName,
		)
		//dsn = dsn + "&loc=Asia%2FShanghai"

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		fmt.Println("数据库连接失败", err)
		return nil
	}
	return db
}
