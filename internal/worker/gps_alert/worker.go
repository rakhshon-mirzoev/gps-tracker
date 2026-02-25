package gps_alert

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"go.uber.org/zap"
)

const (
	// GPSMissingThreshold — сотрудник считается "пропавшим" если не слал GPS дольше этого времени.
	GPSMissingThreshold = 10 * time.Minute
)

// Worker проверяет каждые `interval` сотрудников, не отправлявших GPS
// дольше GPSMissingThreshold, и публикует событие в Kafka.
//
// Уведомление отправляется ровно ОДИН РАЗ пока сотрудник не вернётся:
//   - при обнаружении "пропажи" → SetAlerted (запись в gps_alert_states) → публикуем в Kafka
//   - при получении нового GPS от сотрудника → Resolve (удаление записи)
//   - пока запись существует → IsAlerted == true → пропускаем
type Worker struct {
	locationRepo repository.LocationHistoryRepository
	alertRepo    repository.GpsAlertStatesRepository
	publisher    events.NotificationPublisher
	interval     time.Duration
	stop         chan struct{}
	wg           sync.WaitGroup
}

func New(
	locationRepo repository.LocationHistoryRepository,
	alertRepo repository.GpsAlertStatesRepository,
	publisher events.NotificationPublisher,
	interval time.Duration,
) *Worker {
	return &Worker{
		locationRepo: locationRepo,
		alertRepo:    alertRepo,
		publisher:    publisher,
		interval:     interval,
		stop:         make(chan struct{}),
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.wg.Add(1)
	defer w.wg.Done()

	logger.Log.Info("gps alert worker started", zap.Duration("interval", w.interval))

	for {
		select {
		case <-w.stop:
			logger.Log.Info("gps alert worker stopping")
			return
		case <-ticker.C:
			w.process(context.Background())
		}
	}
}

func (w *Worker) Stop() {
	close(w.stop)
	w.wg.Wait()
}

func (w *Worker) process(ctx context.Context) {
	missing, err := w.locationRepo.GetStaffMissingGPS(ctx, GPSMissingThreshold)
	if err != nil {
		logger.Log.Error("gps alert worker: get missing gps staff", zap.Error(err))
		return
	}

	now := time.Now()

	for _, staff := range missing {
		// Проверяем в БД — уже ли отправляли алерт для этого сотрудника.
		// Пока сотрудник не пришлёт новый GPS (Resolve не вызван) — пропускаем.
		alerted, err := w.alertRepo.IsAlerted(ctx, staff.StaffID)
		if err != nil {
			logger.Log.Error("gps alert worker: check alerted",
				zap.Stringer("staff_id", staff.StaffID),
				zap.Error(err),
			)
			continue
		}
		if alerted {
			continue // уже уведомили, ждём пока сотрудник «вернётся»
		}

		since := now.Sub(staff.CapturedAt).Round(time.Minute)

		event := events.NotificationEvent{
			StaffID:   staff.StaffID,
			GroupID:   staff.GroupID,
			FullName:  staff.FullName,
			Title:     "GPS не получен",
			Body:      fmt.Sprintf("Сотрудник %s не отправлял GPS уже %s", staff.FullName, since),
			Type:      "gps_missing",
			CreatedAt: now,
		}

		// Сохраняем в БД ДО публикации — если Kafka упадёт, не потеряем состояние
		if err := w.alertRepo.SetAlerted(ctx, staff.StaffID); err != nil {
			logger.Log.Error("gps alert worker: set alerted",
				zap.Stringer("staff_id", staff.StaffID),
				zap.Error(err),
			)
			continue
		}

		if err := w.publisher.PublishNotification(ctx, event); err != nil {
			logger.Log.Error("gps alert worker: publish notification",
				zap.Stringer("staff_id", staff.StaffID),
				zap.Error(err),
			)
			// Не откатываем SetAlerted — при следующем запуске всё равно не будем слать повторно.
			// Kafka-сообщение можно будет поймать через retry на уровне producer.
		}
	}
}
