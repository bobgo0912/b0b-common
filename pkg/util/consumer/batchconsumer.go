package consumer

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type BatchConsumer[T any] struct {
	EventQueue     chan *T
	BatchSize      int
	Workers        int
	LingerTime     time.Duration
	BatchProcessor func(batch []*T) error
	ErrHandler     func(err error, batch []*T)
}

func NewBatchConsumer[T any]() *BatchConsumer[T] {
	return &BatchConsumer[T]{
		EventQueue: make(chan *T, 10000),
		BatchSize:  100,
		Workers:    10,
		LingerTime: 100 * time.Millisecond,
		BatchProcessor: func(batch []*T) error {
			return nil
		},
		ErrHandler: func(err error, batch []*T) {
			zap.S().Error("consumer type=", batch, " err=", err.Error())
		},
	}
}

func (c *BatchConsumer[T]) Start(ctx context.Context) {
	for i := 0; i < c.Workers; i++ {
		go func() {
			batch := make([]*T, 0)
			lingerTimer := time.NewTimer(0)
			if !lingerTimer.Stop() {
				<-lingerTimer.C
			}
			defer lingerTimer.Stop()

			for {
				select {
				case msg := <-c.EventQueue:
					batch = append(batch, msg)
					if len(batch) != c.BatchSize {
						if len(batch) == 1 {
							lingerTimer.Reset(c.LingerTime)
						}
						break
					}
					if err := c.BatchProcessor(batch); err != nil {
						c.ErrHandler(err, batch)
					}
					if !lingerTimer.Stop() {
						<-lingerTimer.C
					}
					batch = make([]*T, 0)
				case <-lingerTimer.C:
					if err := c.BatchProcessor(batch); err != nil {
						c.ErrHandler(err, batch)
					}
					batch = make([]*T, 0)
				case <-ctx.Done():
					zap.S().Info("MongoOrderUpdateEvent done")
					return
				}
			}
		}()
	}

}

// options
func (c *BatchConsumer[T]) name() {

}
