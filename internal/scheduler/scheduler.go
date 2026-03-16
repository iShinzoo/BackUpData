package scheduler

import (
	"github.com/iShinzoo/BackUpData/pkg/logger"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {

	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron: c,
	}
}

func (s *Scheduler) AddJob(schedule string, job func()) error {

	_, err := s.cron.AddFunc(schedule, job)

	return err
}

func (s *Scheduler) Start() {

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}

	logg.Info("Schedular Started...")

	s.cron.Start()
}

func (s *Scheduler) Stop() {

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}

	logg.Info("Schedular Stopping...")

	s.cron.Stop()

}
