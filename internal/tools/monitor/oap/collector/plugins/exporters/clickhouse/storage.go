package clickhouse

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/panjf2000/ants"
	"github.com/spf13/cast"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda/internal/tools/monitor/oap/collector/lib"
)

type BatchBuilder interface {
	BuildBatch(ctx context.Context, sourceBatch interface{}) ([]driver.Batch, error)
}

type StorageConfig struct {
	CurrencyNum int `file:"currency_num" default:"100"`
	RetryNum    int `file:"retry_num" default:"5"`
}

type Storage struct {
	logger     logs.Logger
	cfg        *StorageConfig
	batchCh    chan interface{}
	sqlBuilder BatchBuilder

	cancel context.CancelFunc
	ctx    context.Context

	workerWg  sync.WaitGroup
	sendPool  *ants.Pool
	closeOnce sync.Once

	mu     sync.RWMutex
	closed bool
}

func (st *Storage) Start(ctx context.Context) error {
	builders := st.cfg.CurrencyNum
	if builders < 1 {
		builders = 1
	}
	available := lib.AvailableCPUs()
	if proc := os.Getenv("DEBUG_MAX_PROC"); proc != "" {
		available = cast.ToInt(proc)
	}
	if builders < available {
		builders = available
	}

	ctx, cancel := context.WithCancel(ctx)
	st.ctx = ctx
	st.cancel = cancel

	st.logger.Infof("ants new pool, %d", st.cfg.CurrencyNum)
	pool, err := ants.NewPool(st.cfg.CurrencyNum, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	st.sendPool = pool
	st.batchCh = make(chan interface{}, available)
	st.logger.Infof("batch channel, %d", available)

	for i := 0; i < builders; i++ {
		st.workerWg.Add(1)
		go st.handleBatch()
	}
	return nil
}

func (st *Storage) WriteBatchAsync(batch interface{}) {
	for {
		st.mu.RLock()
		if st.closed {
			st.mu.RUnlock()
			st.logger.Warnf("storage closed, dropping batch")
			return
		}
		st.mu.RUnlock()

		select {
		case st.batchCh <- batch:
			return
		case <-st.ctx.Done():
			st.logger.Warnf("storage shutting down, dropping batch")
			return
		}
	}
}

func (st *Storage) handleBatch() {
	defer st.workerWg.Done()

	for {
		select {
		case <-st.ctx.Done():
			return
		case items, ok := <-st.batchCh:
			if !ok {
				return
			}
			batches, err := st.sqlBuilder.BuildBatch(st.ctx, items)
			if err != nil {
				st.logger.Errorf("construct batch: %s", err)
				continue
			}
			for i := range batches {
				st.dispatchSend(batches[i])
			}
		}
	}
}

func (st *Storage) dispatchSend(b driver.Batch) {
	if err := st.sendPool.Submit(func() {
		st.logger.Info("submitted batch job")
		st.sendBatch(b)
	}); err != nil {
		st.logger.Errorf("ck_send_trace submit_error err=%s", err)
		st.abortBatch(b)
	}
}

func (st *Storage) sendBatch(b driver.Batch) {
	start := time.Now()
	backoffDelay := time.Second
	maxBackoffDelay := 30 * time.Second

	for i := 0; i < st.cfg.RetryNum; i++ {
		select {
		case <-st.ctx.Done():
			st.logger.Warnf("ck_send_trace ctx_cancel attempt=%d rows=%d dur=%s", i, b.Rows(), time.Since(start))
			st.abortBatch(b)
			return
		default:
		}

		if b.IsSent() {
			return
		}
		if err := b.Send(); err != nil {
			st.logger.Warnf("ck_send_trace fail attempt=%d rows=%d err=%s backoff=%s", i, b.Rows(), err, backoffDelay)
			if !st.waitWithContext(backoffDelay + time.Duration(rand.Intn(int(backoffDelay/2)))) {
				st.logger.Warnf("ck_send_trace backoff_ctx_cancel attempt=%d rows=%d dur=%s", i, b.Rows(), time.Since(start))
				st.abortBatch(b)
				return
			}
			backoffDelay *= 2
			if backoffDelay > maxBackoffDelay {
				backoffDelay = maxBackoffDelay
			}
			continue
		} else {
			// Only log slow successes to avoid noise.
			if dur := time.Since(start); dur != 0 {
				st.logger.Infof("ck_send_trace done rows=%d dur=%s attempt=%d", b.Rows(), dur, i)
			}
			return
		}
	}
	st.logger.Warnf("ck_send_trace exhausted rows=%d dur=%s attempts=%d", b.Rows(), time.Since(start), st.cfg.RetryNum)
	st.abortBatch(b)
}

func (st *Storage) Close() error {
	st.closeOnce.Do(func() {
		st.mu.Lock()
		st.closed = true
		close(st.batchCh)
		st.mu.Unlock()

		st.workerWg.Wait()
		if st.sendPool != nil {
			st.sendPool.Release()
		}

		if st.cancel != nil {
			st.cancel()
		}
	})
	return nil
}

func (st *Storage) waitWithContext(d time.Duration) bool {
	select {
	case <-st.ctx.Done():
		return false
	case <-time.After(d):
		return true
	}
}

func (st *Storage) abortBatch(b driver.Batch) {
	if err := b.Abort(); err != nil {
		st.logger.Errorf("abort batch: %s", err)
	}
}
