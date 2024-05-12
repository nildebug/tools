package utils

import (
	"context"
	"fmt"
	"github.com/nildebug/tools/log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ListenExitChan
//
//	@Description: 监听程序退出信号
func ListenExitChan() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

var LocalIPs []string

// GetLocalIPs
//
//	@Description: 获取本地IP地址 过滤私有地址
func GetLocalIPs() {

	interfaces, err := net.Interfaces()
	if err != nil {
		log.ZError(context.TODO(), "获取本地绑定IP地址失败", err)
		return
	}

	for _, iface := range interfaces {
		//只取运行状态的接口
		if iface.Flags&net.FlagRunning == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.ZError(context.TODO(), "获取本地绑定IP地址失败", err)
			continue
		}

		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if !ok || ip.IP.IsLoopback() {
				continue
			}
			if ip.IP.To4() != nil && !ip.IP.IsPrivate() {
				LocalIPs = append(LocalIPs, ip.IP.String())
			}
		}
	}

	log.ZDebug(context.TODO(), "localIP", "len", len(LocalIPs), "ips", LocalIPs)
}

// GetRandLocalIP
//
//	@Description: 获取本地随机IP
//	@return string
func GetRandLocalIP() string {
	if len(LocalIPs) == 0 {
		return ""
	}
	index := rand.Intn(len(LocalIPs))
	return LocalIPs[index]
}

// GetCostTime
//
//	@Description: 获取耗时 保留最后两位小数
//	@param startTime
//	@return string
func GetCostTime(startTime time.Time) string {
	duration := time.Since(startTime)
	if duration.Seconds() >= 1 {
		return fmt.Sprintf("%.2f(s)", duration.Seconds())
	}
	return fmt.Sprintf("%.2f(ms)", float64(duration.Nanoseconds())/1e6)
}
