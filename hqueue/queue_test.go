package hqueue

import (
    "sync"
    "log"
    "fmt"
    "testing"
)

/**
 @Description: 相关的常量
 */
const (
    CONSUME_COROUTINE_LIMIT = 10
)

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
        taskQueue.MsgQ = NewQueueAndRun(TaskHandler, 5, QUEUE_LIMIT_DEFAULT)
        log.Println(fmt.Sprintf("init create task queue, task coroutine: %d", 5))
    })

    taskQ := TaskQueueInstance()
    if taskQ == nil {
        panic("init create task queue failed.")
    }
}

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

type Task struct {
    Index           int
}

func (t *Task) Done()  {
    fmt.Println(fmt.Sprintf("hahah======: %d", t.Index))
}

func TestQueue_Push(t *testing.T) {
    fmt.Println("begin.....")
    for i := 0; i < 100; i++ {
        TaskQueueInstance().MsgQ.Push(&Task{Index: i})
    }

    fmt.Println("push over, begin wait....")
    TaskQueueInstance().MsgQ.Wait()
    fmt.Println("end......")
}