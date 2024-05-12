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

}
