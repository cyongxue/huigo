package hqueue

import (
	"sync"
)

/**
 @Description: 队列成员item，interface的形式
 */
type TaskWorker interface {
	Done()                          // 加入任务的类型必须具备Done接口
}

/**
 @Description: queue中注册的handler方法
 */
func TaskHandler(item interface{})  {
	taskItem := item.(TaskWorker)

	taskItem.Done()

	return
}

/**
 @Description: 队列queue
 */
type TaskQueue struct {
	MsgQ            *Queue
}

/**
 @Description: 对外通过获取的接口
 */
func TaskQueueInstance() *TaskQueue {
	return taskQueue
}

/**
 @Description: 全局变量
 */
var taskQueue *TaskQueue
var taskQueueOnce sync.Once

/**
 @Description: 在main函数中采用_ import的方式调用该init接口实现初始化
 */
func init()  {
	taskQueueOnce.Do(func() {
		taskQueue = &TaskQueue{}
	})

	taskQ := TaskQueueInstance()
	if taskQ == nil {
		panic("init create task queue failed.")
	}
}