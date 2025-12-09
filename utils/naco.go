package utils

import (
	"MiniPrograms/responsity/conf"
	"bytes"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

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
	dsn := os.Getenv("NACOSDSN")
	if dsn == "" {
		log.Printf("环境变量NACOSDSN未设置")

		//尝试从.env中获取
		err := gotenv.Load()
		if err != nil {
			log.Printf("加载.env文件失败")
		}
		dsn = os.Getenv("NACOSDSN")
	}

	parts := strings.SplitN(dsn, "?", 2)
	host := parts[0]
	params := url.Values{}

	if len(parts) == 2 {
		params, _ = url.ParseQuery(parts[1])
	}

	hostParts := strings.Split(host, ":")
	server = hostParts[0]

	if len(hostParts) > 1 {
		p, _ := strconv.Atoi(hostParts[1])
		port = uint64(p)
	} else {
		port = 8848
	}

	ns = params.Get("namespace")
	if ns == "" {
		ns = "public"
	}

	user = params.Get("username")
	pass = params.Get("password")
	group = params.Get("group")
	dataId = params.Get("dataId")

	return
}
