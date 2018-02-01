package hqueue

import (
    "fmt"
    "testing"
)

/**
 @Description: 相关的常量
 */
const (
    CONSUME_COROUTINE_LIMIT = 10
)

type Task struct {
    Index           int
}

func (t *Task) Done()  {
    fmt.Println(fmt.Sprintf("hahah======: %d", t.Index))
}

func TestQueue_Push(t *testing.T) {

    TaskQueueInstance().MsgQ = NewQueueAndRun(TaskHandler, 5, QUEUE_LIMIT_DEFAULT)

    fmt.Println("begin.....")
    for i := 0; i < 100; i++ {
        TaskQueueInstance().MsgQ.Push(&Task{Index: i})
    }

    fmt.Println("push over, begin wait....")
    TaskQueueInstance().MsgQ.Wait()
    fmt.Println("end......")
}
