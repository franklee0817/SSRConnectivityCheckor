package server

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

var (
	checkTimes     = 5
	defaultMax     = 9999
	checkTimeOut   = 500 * time.Millisecond
	processorLimit = runtime.NumCPU() * 2
	maxResultQueue = 0
)

// Checkor 服务端配置结构体
type Checkor struct {
	Target Server
}

// CheckServers 根据配置信息检查服务连通性
func CheckServers(servers []Server) ConnectivityList {
	fmt.Println("开始检查服务连通性")
	maxResultQueue = len(servers)
	processor := make(chan Checkor, processorLimit)
	resChan := make(chan Connectivity, maxResultQueue)
	finishedChan := make(chan bool)
	defer close(resChan)
	go initProcessor(processor, servers)

	// 开启processorLimit个协程处理
	go process(processor, resChan, finishedChan)

	finished := <-finishedChan
	if !finished {
		panic("write false to finish channel is not allowed")
	}

	resList := make(ConnectivityList, maxResultQueue)
	for i := 0; i < maxResultQueue; i++ {
		fmt.Print(".")
		resList[i] = <-resChan
	}
	fmt.Println("")
	fmt.Println("检查完毕")
	sort.Sort(resList)

	return resList
}

// initProcessor 初始化进程
func initProcessor(processor chan Checkor, servers []Server) {
	for _, server := range servers {
		fmt.Print(".")
		processor <- Checkor{server}
	}
	close(processor)
}

// process 检查服务连接
func process(processor chan Checkor, resChan chan Connectivity, finishedChan chan bool) {
	procCnt := 0
	for {
		checkor, alive := <-processor
		if !alive {
			fmt.Println("")
			finishedChan <- true
		}
		for {
			if procCnt < processorLimit {
				procCnt++
				go checkServerInRoutine(checkor, resChan, &procCnt)
				break
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}

	}
}

func checkServerInRoutine(checkor Checkor, resChan chan Connectivity, routineCnt *int) {
	totalTime := 0
	validDialCnt := 0
	for i := 0; i < checkTimes; i++ {
		cost, err := checkor.Target.Dial(checkTimeOut)
		if err == nil {
			validDialCnt++
			totalTime += cost
		}
	}
	avgCost := 0
	if validDialCnt > 0 {
		avgCost = int(totalTime / validDialCnt)
	} else {
		avgCost = defaultMax
	}

	connRes := Connectivity{
		fmt.Sprintf("[%s] %s (%s)", checkor.Target.Group, checkor.Target.Remarks, checkor.Target.Server),
		avgCost,
	}

	resChan <- connRes
	*routineCnt--
}
