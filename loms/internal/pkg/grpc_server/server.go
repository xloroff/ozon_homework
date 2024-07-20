package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/grpc_server/interceptor"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// AppServer интерфейс для сервера.
type AppServer interface {
	Start() error
	Stop() error
	RegisterAPI(h []APIHandler)
}

type server struct {
	ctx        context.Context
	logger     logger.ILog
	settings   *config.ApplicationParameters
	grpcServer *grpc.Server
}

// NewServer создает экземпляр grpc сервера.
func NewServer(ctx context.Context, l logger.ILog, cnfg *config.ApplicationParameters) AppServer {
	i := interceptor.NewInterceptor(ctx, l)

	grpcServer := grpc.NewServer(

		grpc.ChainUnaryInterceptor(
			i.Logger(),
			i.Panic(),
			i.Validate(),
		),
	)

	reflection.Register(grpcServer)

	return &server{
		ctx:        ctx,
		logger:     l,
		settings:   cnfg,
		grpcServer: grpcServer,
	}
}

// Start запуск сервера.
func (s *server) Start() error {
	s.logger.Warn(s.ctx, "Запуск GRPC сервера...")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.settings.Port))
	if err != nil {
		return fmt.Errorf("Ошибка использования порта %d: %w", s.settings.Port, err)
	}

	go func() {
		err = s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.Errorf(s.ctx, "Ошибка запуска сервера - %w", err)
		}
	}()

	s.logger.Warn(s.ctx, "GRPC cервер запущен...")

	return nil
}

// Stop остановка сервера.
func (s *server) Stop() error {
	defer func() {
		s.logger.Warn(s.ctx, "Остановка GRPC-сервера произведена...")
	}()

	s.grpcServer.GracefulStop()

	return nil
}
