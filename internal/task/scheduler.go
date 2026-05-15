package task

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron *cron.Cron
	log  *zap.Logger
}

func NewScheduler(log *zap.Logger) *Scheduler {
	return &Scheduler{
		cron: cron.New(),
		log:  log,
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	s.log.Info("task scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.log.Info("task scheduler stopped")
}

// AddJob adds a cron job
func (s *Scheduler) AddJob(spec string, name string, job func()) (cron.EntryID, error) {
	id, err := s.cron.AddFunc(spec, func() {
		s.log.Debug("running scheduled task", zap.String("task", name))
		job()
	})
	if err != nil {
		return 0, err
	}
	s.log.Info("scheduled task added", zap.String("task", name), zap.String("spec", spec))
	return id, nil
}

// DeviceStatusCheck creates a periodic device status check task
func (s *Scheduler) DeviceStatusCheck(interval time.Duration, checkFunc func()) {
	spec := "0 */5 * * * *" // every 5 minutes
	s.AddJob(spec, "device-status-check", func() {
		checkFunc()
	})
}

// AlarmCleanup creates a daily alarm cleanup task
func (s *Scheduler) AlarmCleanup(hour int, cleanupFunc func()) {
	spec := "0 0 " + formatHour(hour) + " * * *"
	s.AddJob(spec, "alarm-cleanup", func() {
		cleanupFunc()
	})
}

func formatHour(i int) string {
	if i < 10 {
		return "0" + itoa(i)
	}
	return itoa(i)
}

func itoa(i int) string {
	return time.Duration(i).String()
}
