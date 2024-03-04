package fixer

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/YanLong-L/migratiox/logger"
	"github.com/YanLong-L/migratiox/migrator/events"
	"github.com/YanLong-L/migratiox/migrator/fixer"
	"github.com/YanLong-L/migratiox/saramax"
	"gorm.io/gorm"
	"time"
)

type Consumer struct {
	client   sarama.Client
	l        logger.Logger
	srcFirst *fixer.OverrideFixer
	dstFirst *fixer.OverrideFixer
	topic    string
}

func NewConsumer(
	client sarama.Client,
	l logger.Logger,
	topic string,
	src *gorm.DB,
	dst *gorm.DB) (*Consumer, error) {
	srcFirst, err := fixer.NewOverrideFixer(src, dst)
	if err != nil {
		return nil, err
	}
	dstFirst, err := fixer.NewOverrideFixer(dst, src)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		client:   client,
		l:        l,
		srcFirst: srcFirst,
		dstFirst: dstFirst,
		topic:    topic,
	}, nil
}

// Start 这边就是自己启动 goroutine 了
func (r *Consumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator-fix",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{r.topic},
			saramax.NewHandler[events.InconsistentEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *Consumer) Consume(msg *sarama.ConsumerMessage, t events.InconsistentEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch t.Direction {
	case "SRC":
		return r.srcFirst.Fix(ctx, t.ID)
	case "DST":
		return r.dstFirst.Fix(ctx, t.ID)
	}
	return errors.New("未知的校验方向")
}
