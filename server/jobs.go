package server

import (
	"context"
	"sync"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/google/logger"
)

type JobCache struct {
	Res *Response
	// Update chan
	sync.RWMutex
}

var devicesCache = &JobCache{}

type DevicesJob struct {
	ctx context.Context
}

func StartJobs(ctx context.Context) {
	jobrunner.Start()
	job := DevicesJob{ctx: ctx}
	err := jobrunner.Schedule("@every 1m", job)
	if err != nil {
		logger.Errorln(err)
	}
	job.Run()
}

func (e DevicesJob) Run() {
	ctx, cancel := context.WithTimeout(e.ctx, time.Minute)
	defer cancel()

	logger.Infoln("Run devices...")

	devices, err := devices(ctx)
	if err != nil {
		logger.Errorln("Job failed with", err)
		return
	}

	devicesCache.Lock()
	defer devicesCache.Unlock()

	devicesCache.Res = devices
}

func GetDevices() *Response {
	rLock := devicesCache.RLocker()
	rLock.Lock()
	defer rLock.Unlock()
	return devicesCache.Res
}
