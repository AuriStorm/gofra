package utils

import (
	"flag"
)

type CliArgs struct {
	Port              *int
	DefaultTimeoutSec *int
	MaxQueueCnt       *int
	QueueSize         *int
}

func MustReadArgs() CliArgs {
	portPtr := flag.Int("port", 0, "set port for app server")
	defaultTimeoutPtr := flag.Int("default-timeout-sec", 0, "set timeout for registered routes")
	maxQueueCntPtr := flag.Int("max-queues", 0, "set max cueues for whole service")
	queueSizePtr := flag.Int("queue-size", 0, "set maximum size per queue")

	flag.Parse()

	return CliArgs{portPtr, defaultTimeoutPtr, maxQueueCntPtr, queueSizePtr}
}
