package server

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ParseSubscribe 从订阅的url中获取ssr内容
func ParseSubscribe(url string) []Server {
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
		server := buildServer(decodeNode)

		servers = append(servers, server)
	}
	fmt.Println("解读订阅信息完毕")
	return servers
}

func buildServer(decodeNode string) Server {
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
