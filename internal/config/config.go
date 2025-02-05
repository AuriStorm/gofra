package config

import (
	"fmt"
	"gofra/internal/utils"
	"log"
	"time"
)

const (
	defaultPort               int32         = 8080
	defaultRouteTimeoutSec    int           = 10
	defaultShutdownTimeoutSec time.Duration = 15
	defaultMaxQueueCnt        int32         = 4
	defaultQueueSize          int32         = 10
)

type AppConfig struct {
	Addr                   string
	RouteDefaultTimeoutSec int
	ShutdownTimeoutSec     time.Duration
}

func (ac *AppConfig) MustLoad(cliArgs utils.CliArgs) *AppConfig {
	var port int32
	var routeDefaultTimeoutSec int

	if cliArgs.Port == nil || *cliArgs.Port == 0 {
		port = defaultPort
	} else {
		port = int32(*cliArgs.Port)
	}

	if cliArgs.DefaultTimeoutSec == nil || *cliArgs.DefaultTimeoutSec == 0 {
		routeDefaultTimeoutSec = defaultRouteTimeoutSec
	} else {
		routeDefaultTimeoutSec = int(*cliArgs.DefaultTimeoutSec)
	}

	ac.Addr = fmt.Sprintf(":%d", port)
	ac.RouteDefaultTimeoutSec = routeDefaultTimeoutSec
	ac.ShutdownTimeoutSec = defaultShutdownTimeoutSec

	log.Println("app config applied successfully")
	return ac
}

type StorageConfig struct {
	MaxQueueCnt int32
	QueueSize   int32
}

func (sc *StorageConfig) MustLoad(cliArgs utils.CliArgs) *StorageConfig {
	var maxQCnt int32
	var qSize int32

	if cliArgs.MaxQueueCnt == nil || *cliArgs.MaxQueueCnt == 0 {
		maxQCnt = defaultMaxQueueCnt
	} else {
		maxQCnt = int32(*cliArgs.MaxQueueCnt)
	}

	if cliArgs.QueueSize == nil || *cliArgs.QueueSize == 0 {
		qSize = defaultQueueSize
	} else {
		qSize = int32(*cliArgs.QueueSize)
	}

	sc.MaxQueueCnt = maxQCnt
	sc.QueueSize = qSize

	log.Println("storage config applied successfully")
	return sc
}
