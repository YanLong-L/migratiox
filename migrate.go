package migratio

import (
	"github.com/YanLong-L/migratiox/gormx/connpool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 做数据迁移的前期配置
type MConfig struct {
	BaseDSN    string   // 源库的连接信息
	TargetDSN  string   // 目标库的连接信息
	KafkaAddrs []string // kafka连接信息
	Tables     []string // key: 源库的表名，迁移到目标表的表名
	AdminPort  int64    // 管理各阶段的gin sever 的端口地址
	MPattern   string
}

func NewConfig(baseDSN string,
	targetDSN string,
	kafkaAddrs []string,
	tables []string,
	adminPort int64,
	MPattern string,
) *MConfig {
	var p = MPattern
	if MPattern == "" {
		p = "SRC_ONLY"
	}
	return &MConfig{BaseDSN: baseDSN,
		TargetDSN:  targetDSN,
		KafkaAddrs: kafkaAddrs,
		Tables:     tables,
		AdminPort:  adminPort,
		MPattern:   p,
	}
}

type SrcDB *gorm.DB
type DstDB *gorm.DB

// InitSRC 初始化源库
func InitSRC(config MConfig) SrcDB {
	return InitDB(config.BaseDSN)
}

// InitDST 初始化目标库
func InitDST(config MConfig) DstDB {
	return InitDB(config.TargetDSN)
}

// InitDoubleWritePool 初始化gorm的connpool，用于在初始化gorm时配置
func InitDoubleWritePool(config MConfig, src SrcDB, dst DstDB) *connpool.DoubleWritePool {
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
