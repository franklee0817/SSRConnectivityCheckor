package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

// Configs 服务配置
type Configs struct {
	Random         bool     `json:"random"`
	AuthPass       string   `json:"authPass"`
	UserOnlinePac  bool     `json:"userOnlinePac"`
	TTL            int      `json:"TTL"`
	Global         bool     `json:"global"`
	ReconnectTimes int      `json:"reconnectTimes"`
	Index          int      `json:"index"`
	ProxyType      int      `json:"proxyType"`
	ProxyHost      string   `json:"proxyHost"`
	AuthUser       string   `json:"authUser"`
	ProxyAuthPass  string   `json:"proxyAuthPass"`
	IsDefault      bool     `json:"isDefault"`
	PacURL         string   `json:"pacUrl"`
	Servers        []Server `json:"configs"`
}

// Server 服务配置详情
type Server struct {
	Enable        bool   `json:"enable"`
	Password      string `json:"password"`
	Method        string `json:"method"`
	Remarks       string `json:"remarks"`
	Server        string `json:"server"`
	Obfs          string `json:"obfs"`
	Protocol      string `json:"protocol"`
	Group         string `json:"group"`
	ServerPort    int    `json:"server_port"`
	RemarksBase64 string `json:"remarks_base64"`
}

// ConnectivityList 节点连通性数组
type ConnectivityList []Connectivity

// Connectivity 节点连通性
type Connectivity struct {
	Name  string
	Delay int
}

func (cl ConnectivityList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}
func (cl ConnectivityList) Len() int {
	return len(cl)
}
func (cl ConnectivityList) Less(i, j int) bool {
	return cl[i].Delay < cl[j].Delay
}

//Dial 测试服务连通性，通过则返回连接耗时，否则报错
func (server *Server) Dial(timeout time.Duration) (int, error) {
	startTime := time.Now().Nanosecond()
	addr := fmt.Sprintf("%s:%v", server.Server, server.ServerPort)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	endTime := time.Now().Nanosecond()
	cost := int((endTime - startTime) / 1000000)
	if cost <= 0 {
		return 0, errors.New("时间获取失败")
	}

	return cost, err
}

// LoadFileConf 读取文件配置
func (sc *Configs) LoadFileConf(filename string) {
	configStr, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Error occur: %s", err))
	}

	err = json.Unmarshal(configStr, sc)
	if err != nil {
		panic(fmt.Sprintf("Error occur: %s", err))
	}
}
