package migratiox

// MConfig 做数据迁移的前期配置
type MConfig struct {
	BaseDSN    string   // 源库的连接信息
	TargetDSN  string   // 目标库的连接信息
	KafkaAddrs []string // kafka连接信息
	Tables     []string // key: 源库的表名，迁移到目标表的表名
	AdminPort  int64    // 管理各阶段的gin sever 的端口地址
	MPattern   string
}

// NewMConfig 初始化Config
func NewMConfig(baseDSN string,
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

func InitMigrateApp() {

}
