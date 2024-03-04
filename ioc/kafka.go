package ioc

import (
	"github.com/IBM/sarama"
	"github.com/YanLong-L/migratiox"
	"github.com/YanLong-L/migratiox/migrator/events/fixer"
	"github.com/YanLong-L/migratiox/saramax"
)

func InitKafka(config migratiox.MConfig) sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(config.KafkaAddrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConsumers 注册修复数据的Consumer
func NewConsumers(
	fix *fixer.Consumer,
) []saramax.Consumer {
	return []saramax.Consumer{
		fix,
	}
}
