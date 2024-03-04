package ioc

import (
	"github.com/IBM/sarama"
	"github.com/YanLong-L/migratiox/ginx"
	"github.com/YanLong-L/migratiox/gormx/connpool"
	"github.com/YanLong-L/migratiox/logger"
	"github.com/YanLong-L/migratiox/migrator/events"
	"github.com/YanLong-L/migratiox/migrator/events/fixer"
	"github.com/YanLong-L/migratiox/migrator/scheduler"
	"github.com/gin-gonic/gin"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.Logger,
	src SrcDB,
	dst DstDB,
	client sarama.Client) *fixer.Consumer {
	res, err := fixer.NewConsumer(client, l,
		topic, src, dst)
	if err != nil {
		panic(err)
	}
	return res
}

func InitMigradatorProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer(p, topic)
}

func InitMigratorWeb(
	l logger.Logger,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer,
) *ginx.Server {
	// 在这里，有多少张表，你就初始化多少个 scheduler
	intrSch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	engine := gin.Default()

	intrSch.RegisterRoutes(engine.Group("/migrator"))
	//intrSch.RegisterRoutes(engine.Group("/migrator/interactive"))
	addr := viper.GetString("migrator.web.addr")
	return &ginx.Server{
		Addr:   addr,
		Engine: engine,
	}
}
