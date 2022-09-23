package util

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
)

var (
	nacosUtilInstance *Nacos
	nacosUtilOnce     sync.Once
)

type Nacos struct {
	// 配置中心客户端
	configClient config_client.IConfigClient
	// 服务发现客户端
	namingClient naming_client.INamingClient
	// 分组
	group string
}

// Nacos 工具
func NacosUtil() *Nacos {
	nacosUtilOnce.Do(func() {
		nacosUtilInstance = new(Nacos)
	})
	return nacosUtilInstance
}

// 服务注册发现客户端
// @param namingClient 服务注册发现客户端
func (u *Nacos) WithNamingClient(namingClient naming_client.INamingClient) *Nacos {
	u.namingClient = namingClient
	return u
}

// 配置中心客户端
// @param configClient 配置中心客户端
func (u *Nacos) WithConfigClient(configClient config_client.IConfigClient) *Nacos {
	u.configClient = configClient
	return u
}

// 分组
// @params group 分组名
func (u *Nacos) Group(group string) *Nacos {
	u.group = group
	return u
}

// 获取配置
// @param dataId 数据ID
// @param defaultValue 默认值
func (u Nacos) GetConfig(dataId string, defaultValue string) string {
	// 获取配置
	configData, err := u.configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  u.group,
	})

	if err != nil {
		return defaultValue
	}

	return configData
}

// 监听配置
// @param dataId 数据ID
// @param onChangeCallback 配置变更回调函数
func (u Nacos) ListenConfig(dataId string, onChangeCallback func(namespace string, group string, dataId string, data string)) {
	// 监听配置
	u.configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    u.group,
		OnChange: onChangeCallback,
	})
}

// 注册服务实例
func (u Nacos) RegisterInstance(instanceParam vo.RegisterInstanceParam) (bool, error) {
	return u.namingClient.RegisterInstance(instanceParam)
}
