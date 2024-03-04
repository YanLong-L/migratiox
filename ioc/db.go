package ioc

import (
	"github.com/YanLong-L/migratiox"
	_ "github.com/YanLong-L/migratiox"
	"github.com/YanLong-L/migratiox/gormx/connpool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SrcDB *gorm.DB
type DstDB *gorm.DB

// InitSRC 初始化源库
func InitSRC(config migratiox.MConfig) SrcDB {
	return InitDB(config.BaseDSN)
}

// InitDST 初始化目标库
func InitDST(config migratiox.MConfig) DstDB {
	return InitDB(config.TargetDSN)
}

// InitDoubleWritePool 初始化gorm的connpool，用于在初始化gorm时配置
func InitDoubleWritePool(config migratiox.MConfig, src SrcDB, dst DstDB) *connpool.DoubleWritePool {
	return connpool.NewDoubleWritePool(src.ConnPool, dst.ConnPool, config.MPattern)
}

// InitBizDB 这个是进行数据迁移时，业务用的，支持双写的 DB
func InitBizDB(pool *connpool.DoubleWritePool) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: pool,
	}))
	if err != nil {
		panic(err)
	}
	return db
}

// InitDB 初始化一个gormDB 通过一个key,区分加载配置文件中的哪个库的数据库dsn
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
