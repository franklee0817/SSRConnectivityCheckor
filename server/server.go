package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
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

// PullSubscribe 从订阅的url中获取ssr内容
func PullSubscribe(url string) []Server {
	// 从订阅地址读取内容
	fmt.Println("开始获取订阅信息")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取订阅信息失败：", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取订阅信息失败：", err)
	}

	// 解析订阅内容
	context, err := base64Decode(string(body))
	if err != nil {
		fmt.Println("获取订阅信息失败：", err)
	}
	fmt.Println("获取订阅信息完毕，开始解读")

	servers := []Server{}
	for {
		// 开始解析订阅内容
		idx := strings.Index(context, "\n")
		if idx <= 0 {
			break
		}
		// 跳过 前缀ssr://
		encodedNode := context[6:idx]
		context = context[idx+1:]
		decodeNode, err := base64Decode(encodedNode)
		if err != nil {
			fmt.Println("解析订阅信息失败：", err)
		}
		server := subscribeToServer(decodeNode)

		servers = append(servers, server)
	}
	fmt.Println("解读订阅信息完毕")
	return servers
}

// subscribeToServer 解析订阅返回的单个服务器配置
func subscribeToServer(decodeNode string) Server {
	// host
	idx := strings.Index(decodeNode, ":")
	host := decodeNode[:idx]
	decodeNode = decodeNode[idx+1:]

	// port
	idx = strings.Index(decodeNode, ":")
	portStr := decodeNode[:idx]
	decodeNode = decodeNode[idx+1:]
	port, _ := strconv.Atoi(portStr)

	// obfs
	idx = strings.Index(decodeNode, ":")
	protocol := decodeNode[:idx]
	decodeNode = decodeNode[idx+1:]

	// method
	idx = strings.Index(decodeNode, ":")
	method := decodeNode[:idx]
	decodeNode = decodeNode[idx+1:]

	// obfs
	idx = strings.Index(decodeNode, ":")
	obfs := decodeNode[:idx]
	decodeNode = decodeNode[idx+1:]

	// password， 这里暂时用不到password，所以解出来出错也暂不处理
	idx = strings.Index(decodeNode, "/")
	password := decodeNode[:idx]
	decodeNode = decodeNode[idx+2:] // 跳过/?的?
	password, _ = base64Decode(password)

	// obfsparam
	idx = strings.Index(decodeNode, "&")
	obfsparam := decodeNode[:idx]
	obfsparamIdx := strings.Index(obfsparam, "=")
	obfsparam = obfsparam[obfsparamIdx+1:]
	decodeNode = decodeNode[idx+1:]

	// protoparam
	idx = strings.Index(decodeNode, "&")
	protoparam := decodeNode[:idx]
	protoparamIdx := strings.Index(protoparam, "=")
	protoparam = protoparam[protoparamIdx+1:]
	decodeNode = decodeNode[idx+1:]

	// remarks
	idx = strings.Index(decodeNode, "&")
	remarks := decodeNode[:idx]
	remarksIdx := strings.Index(remarks, "=")
	remarksBase64 := remarks[remarksIdx+1:]
	decodeNode = decodeNode[idx+1:]
	remarks, _ = base64Decode(remarksBase64)

	// group
	group := decodeNode
	groupIdx := strings.Index(group, "=")
	group = group[groupIdx+1:]
	group, _ = base64Decode(group)

	server := Server{
		Enable:        true,
		Password:      password,
		Method:        method,
		Remarks:       remarks,
		Server:        host,
		Obfs:          obfs,
		Protocol:      protocol,
		Group:         group,
		ServerPort:    port,
		RemarksBase64: remarksBase64,
	}

	return server
}

func base64Decode(enStr string) (string, error) {
	deBytes, err := base64.RawURLEncoding.DecodeString(enStr)

	return string(deBytes), err
}
