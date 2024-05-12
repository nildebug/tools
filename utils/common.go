package utils

import (
	"context"
	"fmt"
	"github.com/nildebug/tools/log"
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
//	@Description: 获取本地IP地址
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
			if ip.IP.To4() != nil {
				LocalIPs = append(LocalIPs, ip.IP.String())
				log.ZDebug(context.TODO(), "localIP", "ip", ip.IP.String(), "name", iface.Name, "Flags", iface.Flags)
			}
		}
	}
}

// IsPrivateIP
//
//	@Description: 绑定IP地址是否为10 172 192 段
//	@param ip
//	@return bool
func IsPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		// 10.0.0.0/8
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		// 172.16.0.0/12
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		// 192.168.0.0/16
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
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
