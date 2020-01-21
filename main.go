package main

import (
	"SSRConnectivityCheckor/pathloader"
	"SSRConnectivityCheckor/server"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// defaultConfig 配置文件内容结构体，config_file: ssr服务器json配置文件路径，subscribe_url: ssr订阅地址
type defaultConfig struct {
	ConfigFile   string `json:"config_file"`
	SubscribeURL string `json:"subscribe_url"`
}

func main() {
	// 读取运行参数
	var subscribeURL, configFile string
	flag.StringVar(&subscribeURL, "s", "", "SSR订阅地址")
	flag.StringVar(&configFile, "f", "", "SSR配置文件路径")
	flag.Parse()

	if len(subscribeURL) > 0 {
		fromSubscribe(subscribeURL)
	} else if len(configFile) > 0 {
		fromConfigFile(configFile)
	} else {
		// 获取用户Home目录下的配置文件内容
		userHome, _ := pathloader.Home()
		fileName := userHome + "/.ssr_scanner_conf"
		file, _ := os.Open(fileName)
		decoder := json.NewDecoder(file)

		conf := defaultConfig{}
		err := decoder.Decode(&conf)
		if err != nil {
			str := fmt.Sprintf("找不到默认配置文件: %s/.ssr_scanner_conf", userHome)
			fmt.Println(str)
			return
		}
		if len(conf.SubscribeURL) > 0 {
			fromSubscribe(conf.SubscribeURL)
		} else if len(conf.ConfigFile) > 0 {
			fromConfigFile(conf.ConfigFile)
		} else {
			str := fmt.Sprintf("配置文件为空： %s/.ssr_scanner_conf", userHome)
			fmt.Println(str)
			return
		}
	}
}

// fromSubscribe 从订阅地址获取服务器列表并检测连通性
func fromSubscribe(url string) {
	servers := server.PullSubscribe(url)
	cl := server.CheckServers(servers)
	printCl(cl)

	return
}

// fromConfigFile 从配置文件获取服务器列表并检测连通性
func fromConfigFile(configFile string) {
	sc := &server.Configs{}
	sc.LoadFileConf(configFile)
	cl := server.CheckServers(sc.Servers)
	printCl(cl)

	return
}

func printCl(cl server.ConnectivityList) {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("===========================================================================")
	fmt.Println("延迟最低的前20的节点为.")
	i := 0
	for _, v := range cl {
		i++
		str := fmt.Sprintf("%d. %s 延迟: %d ms", i, v.Name, v.Delay)
		fmt.Println(str)
		if i >= 20 {
			break
		}
	}
}
