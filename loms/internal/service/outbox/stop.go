package outboxservice

// Stop останавливает сервис отправки в брокер.
func (s *service) Stop() error {
	defer func() {
		s.logger.Warn(s.ctx, "Остановка сервиса отправки эвентов в брокер произведена...")
	}()

	close(s.sendStatuses)
	s.wg.Wait()

	return nil
}
