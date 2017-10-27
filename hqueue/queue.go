package hqueue

import (
    "runtime"
    "container/list"
    "sync"
    "time"
)

const (
    QUEUE_LIMIT_DEFAULT =  10240
)

/**
 @Description：queue的设计结构
 */
type Queue struct {
    queueBuf        *list.List
    queueLimit      int                         // queue缓冲的最大空间限制

    pushBack    chan interface{}                // 用于向list中加入item
    popFront    chan interface{}                // 用于从表示从list中消耗了一个

    suspend     chan bool                       // 用于通知queue暂停
    suspended   bool                            // 标记暂停与否
    stop        chan bool                       // 通知停止
    stopped     bool                            // 标记停止

    currentGoCount     int
    GoLimit            int                      // 限制的协程数
    Handler            func(interface{})        // 回调接口

    wg                 sync.WaitGroup           // 等待协程，该功能是可选
}

/**
 @Description：创建一个queue
 @Param:
 @Return：
 */
func NewQueue(handler func(interface{}), limit int, qLimit int) *Queue {

    if limit <= 0 {
        limit = 1
    }

    newQ := &Queue{
        queueBuf: list.New(),
        queueLimit: qLimit,

        pushBack: make(chan interface{}),
        popFront: make(chan interface{}),
        suspend: make(chan bool),
        suspended: false,
        stop: make(chan bool),
        stopped: false,

        currentGoCount: 0,
        GoLimit: limit,
        Handler: handler,

        wg: sync.WaitGroup{},
    }
    runtime.SetFinalizer(newQ, (*Queue).Stop)
    return newQ
}

/**
@Description：创建一个queue，并开始run
@Param:
@Return：
*/
func NewQueueAndRun(handler func(interface{}), limit int, qLimit int) *Queue {

    if limit <= 0 {
        limit = 1
    }

    newQ := &Queue{
        queueBuf: list.New(),
        queueLimit: qLimit,

        pushBack: make(chan interface{}),
        popFront: make(chan interface{}),
        suspend: make(chan bool),
        suspended: false,
        stop: make(chan bool),
        stopped: false,

        GoLimit: limit,
        Handler: handler,

        wg: sync.WaitGroup{},
    }

    go newQ.run()                   // 创建一个新协程，作为消费协程
    // 参考：http://wiki.jikexueyuan.com/project/the-way-to-go/10.8.html
    // 自定义注册内存回收
    runtime.SetFinalizer(newQ, (*Queue).Stop)
    return newQ
}

/**
 @Description: 开始run
 @Param:
 @Return：
*/
func (q *Queue) Run() {
    go q.run()
}

/**
 @Description: run的实际实现逻辑
 @Param:
 @Return：
*/
func (q *Queue) run() {
    // 退出run，则销毁buf和wait
    defer func() {
        q.wg.Add(-q.queueBuf.Len())
        q.queueBuf = nil
    }()

    // run中采用loop的方式
    for {
        // event
        select {
        case item := <-q.pushBack:
            if q.stopped != true {
                q.queueBuf.PushBack(item)
                q.wg.Add(1)                  // item数作为wait的对象，保证只有消费完，才结
            } else {
                item = nil
            }

        case <-q.popFront:
            q.currentGoCount--                 // 协程结束触发-1，标记协程数少了

        case suspendFlag := <-q.suspend:
            if q.suspended != suspendFlag {
                q.suspended = suspendFlag
                if q.suspended {
                    q.wg.Add(1)          // 暂停，则增加一个event，保证wait一直
                } else {
                    q.wg.Done()
                }
            }

        case <-q.stop:
            q.stopped = true
        }

        // 标记需要停止，且item消费完，则停止
        if q.stopped && (q.queueBuf.Len() <= 0) {
            return
        }

        // 一直消费
        for (q.currentGoCount < q.GoLimit) && (q.queueBuf.Len() > 0) {

            item := q.queueBuf.Front()
            if item != nil {
                taskItem := q.queueBuf.Remove(item)

                q.currentGoCount++               // 即将要创建协程，所以++
                go func() {
                    defer func() {
                        q.popFront <- struct{}{}
                        q.wg.Done()
                    }()

                    q.Handler(taskItem)
                }()
            }
        }
    }
}

/**
 queue停止服务
 */
func (q *Queue) Stop()  {
    q.stop <- true
    // 注销内存回收处理
    runtime.SetFinalizer(q, nil)
}

/**
 加入一个元素
 */
func (q *Queue) Push(item interface{})  {
    for {
        if q.queueBuf.Len() < q.queueLimit {
            // 空间够，则直接push
            q.pushBack <- item
            return
        }

        // 空间不够，则停顿一会
        time.Sleep(2 * time.Second)
    }
    return
}

/**
 等待，stop之后可以采用wait等待消耗完毕
 */
func (q *Queue) Wait() {
    q.wg.Wait()
}

func (q *Queue) Len() (int, int) {
    return q.currentGoCount, q.queueBuf.Len()
}