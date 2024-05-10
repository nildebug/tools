package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// ListenExitChan
//
//	@Description: 监听程序退出信号
func ListenExitChan() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
