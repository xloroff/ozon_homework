package closer

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// Closer абстракция завершения приложения.
type Closer interface {
	Add(f ...func() error)
	Wait()
	Signal()
	CloseAll()
}

type closer struct {
	sync.Mutex
	logger   logger.ILog
	once     sync.Once
	done     chan struct{}
	funcs    []func() error
	shutdown chan os.Signal
}

// NewCloser ожидает сигналы для завершения функций.
func NewCloser(l logger.ILog, sig ...os.Signal) Closer {
	c := &closer{
		logger:   l,
		done:     make(chan struct{}),
		shutdown: make(chan os.Signal, 1),
	}

	if len(sig) > 0 {
		go func() {
			signal.Notify(c.shutdown, sig...)
			<-c.shutdown
			signal.Stop(c.shutdown)
			c.CloseAll()
		}()
	}

	return c
}

// Add добавляет функцию в список завершающих функций.
func (c *closer) Add(f ...func() error) {
	c.Lock()
	c.funcs = append(c.funcs, f...)
	c.Unlock()
}

// Wait подвешиваем ожидание, чтобы ничего не закрылось заранее.
func (c *closer) Wait() {
	<-c.done
}

// Signal отправка сигнала завершения.
func (c *closer) Signal() {
	close(c.shutdown)
}

// CloseAll запуск завершения всех функций.
func (c *closer) CloseAll() {
	c.once.Do(func() {
		ctx := context.Background()
		c.logger.Info(ctx, "Поступил сигнал завершения приложения...")
		defer c.logger.Info(ctx, "Завершение приложения окончено...")

		defer close(c.done)

		c.Lock()
		funcs := make([]func() error, len(c.funcs))
		copy(funcs, c.funcs)
		c.Unlock()

		for i := len(funcs) - 1; i >= 0; i-- {
			err := c.funcs[i]()
			if err != nil {
				c.logger.Errorf(ctx, "Ошибка правильного завершения функции - %v", err)
			}
		}
	})
}
