package main

import (
	"SSRConnectivityCheckor/path_loader"
	"SSRConnectivityCheckor/server"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type defaultConfig struct {
	ConfigFile   string `json:"config_file"`
	SubscribeURL string `json:"subscribe_url"`
}

func main() {
	var subscribeURL, configFile string
	flag.StringVar(&subscribeURL, "s", "", "SSR订阅地址")
	flag.StringVar(&configFile, "f", "", "SSR配置文件路径")
	flag.Parse()

	if len(subscribeURL) > 0 {
		fromSubscribe(subscribeURL)
	} else if len(configFile) > 0 {
		fromConfigFile(configFile)
	} else {
		userHome, _ := path_loader.Home()
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

func fromSubscribe(url string) {
	servers := server.ParseSubscribe(url)
	cl := server.CheckServers(servers)
	printCl(cl)

	return
}

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
