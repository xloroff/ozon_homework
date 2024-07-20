package errgroup

import (
	"context"
	"strings"
	"sync"
)

// Errgroup - структура для управления группой горутин с поддержкой контекста и отмены.
type Errgroup struct {
	ctx           context.Context
	cancel        func()
	funcs         []func() error
	workersQueue  chan struct{}
	mu            sync.RWMutex
	errors        []error
	cancelIfFirst bool
	wg            sync.WaitGroup
}

// errgroupOption настройки инициализации.
type errgroupOption func(*Errgroup)

// NewErrGroup создает новый errgroup с заданным контекстом и количеством воркеров.
func NewErrGroup(ctx context.Context, opts ...errgroupOption) (*Errgroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	errGr := &Errgroup{
		ctx:           ctx,
		cancel:        cancel,
		workersQueue:  make(chan struct{}, 1),
		cancelIfFirst: false,
	}

	for _, option := range opts {
		option(errGr)
	}

	return errGr, ctx
}

// WithCancelFirst опция позволяющая стопать выполнение запросов при появлении первой ошибки.
func WithCancelFirst() func(*Errgroup) {
	return func(s *Errgroup) {
		s.cancelIfFirst = true
	}
}

// WithWorkersCount опция устанавливающая число одновременных запросов.
func WithWorkersCount(count int) func(*Errgroup) {
	return func(s *Errgroup) {
		s.workersQueue = make(chan struct{}, count)
	}
}

// Go добавляет функцию в группу для выполнения.
func (eg *Errgroup) Go(f func() error) {
	eg.funcs = append(eg.funcs, f)
}

// Wait ожидает завершения всех задач и возвращает первую возникшую ошибку, если она была.
func (eg *Errgroup) Wait() error {
	eg.run()
	eg.wg.Wait()

	return eg.FirstError()
}

// run запускает все функции в горутинах и ожидает их завершения.
func (eg *Errgroup) run() {
	for _, f := range eg.funcs {
		select {
		case <-eg.ctx.Done():
			return
		case eg.workersQueue <- struct{}{}:
		}

		eg.wg.Add(1)

		go eg.runFunc(f)
	}
}

// runFunc запускает функцию в горутине.
func (eg *Errgroup) runFunc(f func() error) {
	defer func() {
		<-eg.workersQueue
		eg.wg.Done()
	}()

	if err := f(); err != nil {
		eg.handleError(err)
	}
}

// handleError обработка ошибки и отмена горутин, если errgroup создавалась как первая ошибка = всем стоп.
func (eg *Errgroup) handleError(err error) {
	eg.mu.Lock()
	defer eg.mu.Unlock()

	eg.errors = append(eg.errors, err)
	if eg.cancelIfFirst {
		eg.cancel()
	}
}

// Errors возвращает все ошибки, возникшие при выполнении функций.
func (eg *Errgroup) Errors() []error {
	eg.mu.RLock()
	defer eg.mu.RUnlock()

	errorsCopy := make([]error, len(eg.errors))
	copy(errorsCopy, eg.errors)

	return errorsCopy
}

// FirstError возвращает первую возникшую ошибку или nil, если ошибок не было.
func (eg *Errgroup) FirstError() error {
	eg.mu.RLock()
	defer eg.mu.RUnlock()

	if len(eg.errors) > 0 {
		return eg.errors[0]
	}

	return nil
}

// ErrsToString объединяет все ошибки в одну строку.
func (eg *Errgroup) ErrsToString() string {
	eg.mu.RLock()
	defer eg.mu.RUnlock()

	if len(eg.errors) == 0 {
		return ""
	}

	errorStrings := make([]string, len(eg.errors))
	for i, err := range eg.errors {
		errorStrings[i] = err.Error()
	}

	return strings.Join(errorStrings, "; ")
}
