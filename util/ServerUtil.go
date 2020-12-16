package util

import (
	"net"
	"sync"
)

var (
	serverUtilInstance *Server
	serverUtilOnce     sync.Once
)

func ServerUtil() *Server {
	serverUtilOnce.Do(func() {
		serverUtilInstance = new(Server)
	})
	return serverUtilInstance
}

// 服务器工具
type Server struct {
}

// 获取服务端IP
func (u Server) GetServerIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}
