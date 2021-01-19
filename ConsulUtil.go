// 参考资料：https://studygolang.com/articles/28241?fr=sidebar
package util

//import (
//	"fmt"
//	"github.com/micro/go-micro/client/selector"
//	"github.com/micro/go-micro/registry"
//	"github.com/micro/go-plugins/registry/consul"
//	"time"
//)
//
//// Consul 注册
//var consulRegister registry.Registry
//
//type ConsulUtil struct {
//}
//
//// 注册服务
//func (consulUtil *ConsulUtil) RegisterServer(registerAddrs []string) registry.Registry {
//	consulRegister = consul.NewRegistry(registry.Addrs(registerAddrs...))
//	return consulRegister
//}
//
//// 获取服务
//func (consulUtil *ConsulUtil) GetServer(serverName string) (registryNode *registry.Node) {
//	// 重试次数
//	var retryCount int
//
//	for {
//		// 获取注册服务
//		servers, err := consulRegister.GetService(serverName)
//
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//
//		// 轮询获取其中一个服务信息
//		next := selector.RoundRobin(servers)
//		node, err := next()
//
//		if err == nil {
//			return node
//		}
//
//		retryCount++
//		fmt.Printf("重试：%d，Error: %s\n", retryCount, err.Error())
//
//		if retryCount >= 5 {
//			fmt.Printf("重试失败，Error: %s\n", err.Error())
//			return
//		}
//
//		// 休眠1秒
//		time.Sleep(time.Second * 1)
//	}
//}
