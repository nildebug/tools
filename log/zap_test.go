package log

import (
	"context"
	"fmt"
	"github.com/nildebug/tools/mcontext"
	"testing"
	"time"
)

func TestInitFromConfig(t *testing.T) {
	//fmt.Println("F:/rabbitmq-demo/log/zap_test.go:18")
	if err := InitFromConfig("test", "modeuleName", 5, true, true, "./logs", 1); err != nil {
		panic(err)
	}
	ctx := context.TODO()
	ctx = mcontext.SetValue(ctx, mcontext.CenterKey, "中心key")

	go func() {
		time.Sleep(2 * time.Second)
		SetLogLevel(LogLevelError)

		time.Sleep(5 * time.Second)
		SetLogLevel(LogLevelDubug)
	}()
	for i := 0; i < 100; i++ {
		ZError(ctx, "ssss", fmt.Errorf("error "))
		time.Sleep(time.Second)
		ZDebug(ctx, "zdebug", "name", "sss")
	}
}
