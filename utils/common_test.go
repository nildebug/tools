package utils

import (
	"context"
	"github.com/nildebug/tools/log"
	"testing"
	"time"
)

func TestGetLocalIPs(t *testing.T) {
	log.InitFromConfig("log", "debug", 5, true, true, "", 3)
	startTime := time.Now()
	GetLocalIPs()
	//time.Sleep(1 * time.Second)
	log.ZDebug(context.TODO(), GetCostTime(startTime))

	smax := 10
	smin := 1
	RandInt(smax, smin)

	var s int64 = 1
	var b int64 = 2
	RandInt(s, b)
}
