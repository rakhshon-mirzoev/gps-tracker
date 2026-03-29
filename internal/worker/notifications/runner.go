package notifications

import (
	"context"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
)

func (w *Worker) Run() {
	w.ticker = time.NewTicker(w.interval)
	defer w.ticker.Stop()

	w.wg.Add(1)
	defer w.wg.Done()

	logger.Log.Info("notification worker started")

	for {
		select {
		case <-w.stop:
			logger.Log.Info("notification worker stopping")
			return
		case <-w.ticker.C:
			w.process(context.Background())
		}
	}
}

func (w *Worker) Stop() {
	logger.Log.Info("stopping notifications worker")
	close(w.stop)
	w.wg.Wait()
	logger.Log.Info("notifications worker stopped")
}
