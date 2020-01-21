// Package multi 带有最大协程数的协程控制类
package multi

import (
	"sync"
	"sync/atomic"
	"time"
)

// INIT 初始化状态
const INIT = 0

// PROCESSING 任务执行中
const PROCESSING = 1

// FINISHED 任务已完成
const FINISHED = 2

// Task 多任务结构体， procMax代表最大同时执行的任务数，procCnt代表当前执行中的任务数
type Task struct {
	wg       *sync.WaitGroup
	procMax  int32
	procCnt  int32
	procceed int32
	total    int32
}

// NewTask 创建新实例
func NewTask(procMax int32) *Task {
	mp := &Task{}
	wg := &sync.WaitGroup{}
	mp.wg = wg
	mp.procMax = procMax
	mp.procCnt = 0
	mp.procceed = 0
	mp.total = 0

	return mp
}

// SetTaskTotal 设置任务总数
func (mp *Task) SetTaskTotal(taskCnt uint32) {
	if taskCnt <= 0 {
		panic("任务总数必须大于0")
	}
	if mp.total > 0 {
		panic("不得重复设置任务总数")
	}
	mp.total = int32(taskCnt)
}

// ForceUpdateTaskTotal 强行更新任务总数
func (mp *Task) ForceUpdateTaskTotal(taskCnt uint32) {
	if taskCnt <= 0 {
		panic("任务总数必须大于0")
	}
	atomic.SwapInt32(&mp.total, int32(taskCnt))
}

// Start 开始新任务，若协程数量达到上限则阻塞等待
func (mp *Task) Start() {
	for {
		procCnt := atomic.AddInt32(&mp.procCnt, 1)
		if procCnt > mp.procMax {
			atomic.AddInt32(&mp.procCnt, -1)
			time.Sleep(time.Millisecond)
		} else {
			break
		}
	}
	mp.wg.Add(1)
}

// Done 任务执行完毕
func (mp *Task) Done() {
	atomic.AddInt32(&mp.procCnt, -1)
	mp.wg.Done()
}

// Wait 阻塞进程，等待任务执行完毕
func (mp *Task) Wait() {
	mp.wg.Wait()
}

// Status 查看任务当前执行状态
func (mp *Task) Status() uint8 {
	if mp.procceed == 0 {
		return INIT
	} else if mp.procceed < mp.total {
		return PROCESSING
	} else {
		return FINISHED
	}
}
