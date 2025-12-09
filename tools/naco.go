package tools

import (
	"MiniPrograms/responsity/conf"
	"bytes"
	"log"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// 从nacos获取配置
func GetConfigFromNacos(c *conf.Config) error {
	server, port, namespace, user, pass, group, dataId := ParseNacosDSN()
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: server,
			Port:   port,
			Scheme: "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace,
		Username:            user,
		Password:            pass,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		CacheDir:            "./",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		return err
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err = v.ReadConfig(bytes.NewBufferString(content)); err != nil {
		log.Printf("解析配置文件失败:%v", err)
		return err
	}

	err = v.Unmarshal(c)
	if err != nil {
		log.Printf("反序列化失败:%v", err)
		return err
	}

	return nil

}

func ParseNacosDSN() (server string, port uint64, ns, user, pass, group, dataId string) {
	//开发环境使用.env文件测试
	err := gotenv.Load()
	if err != nil {
		log.Printf("加载.env文件失败")
	}
	server = os.Getenv("SERVER")
	portInt, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Printf("读取端口号失败")
	}

	port = uint64(portInt)
	ns = os.Getenv("NAMESPACE")
	user = os.Getenv("USERNAME")
	pass = os.Getenv("PASSWORD")
	group = os.Getenv("GROUP")
	dataId = os.Getenv("DATAID")
	return server, port, ns, user, pass, group, dataId
}
