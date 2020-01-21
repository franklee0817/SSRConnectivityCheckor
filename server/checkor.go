package server

import (
	"SSRConnectivityCheckor/multi"
	"fmt"
	"runtime"
	"sort"
	"time"
)

var (
	checkTimes   = 5                                  // 每个节点的检查次数
	defaultMax   = 9999                               // 默认节点延迟，当节点不通时使用此值作为节点延迟
	checkTimeOut = 500 * time.Millisecond             // 检查节点延迟的超时时间，超过500毫秒认为节点不通
	procLimit    = runtime.NumCPU() * 2               // 同时开启的任务协程的最大数量
	multiTask    = multi.NewTask(int32(procLimit))    // 多任务控制类
	resChan      = make(chan Connectivity, procLimit) // 检查结果收集通道
)

// Checkor 服务端配置结构体
type Checkor struct {
	Target Server
}

// CheckServers 根据配置信息检查服务连通性
func CheckServers(servers []Server) ConnectivityList {
	fmt.Println("开始检查服务连通性")

	// 设置任务总数量
	slen := len(servers)
	multiTask.SetTaskTotal(uint32(slen))

	go process(servers)

	resList := gatherCheckResult()
	fmt.Println("")
	fmt.Println("检查完毕")
	sort.Sort(resList)

	return resList
}

// gatherCheckResult 收集处理结果
func gatherCheckResult() ConnectivityList {
	resList := ConnectivityList{}
	for {
		res, open := <-resChan
		if !open {
			break
		}
		resList = append(resList, res)
	}

	return resList
}

// process 检查服务连接
func process(servers []Server) {
	for _, serv := range servers {
		multiTask.Start()
		fmt.Print(".")
		go checkServer(serv)
	}
	multiTask.Wait()
	close(resChan)
}

// checkServer 检查服务连通性
func checkServer(serv Server) {
	totalTime := 0
	validDialCnt := 0
	for i := 0; i < checkTimes; i++ {
		// 总共执行checkTimes次检查
		cost, err := serv.Dial(checkTimeOut)
		if err == nil {
			// 统计有效检查次数和总耗时
			validDialCnt++
			totalTime += cost
		}
	}

	// 计算联通耗时均值，若节点不通则使用defaultMax作为节点延迟
	avgCost := 0
	if validDialCnt > 0 {
		avgCost = int(totalTime / validDialCnt)
	} else {
		avgCost = defaultMax
	}

	connRes := Connectivity{
		fmt.Sprintf("[%s] %s (%s)", serv.Group, serv.Remarks, serv.Server),
		avgCost,
	}

	resChan <- connRes
	multiTask.Done()
}
