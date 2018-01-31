// @Description：爬虫相关的基础结构
//               request的封装
package crawler

import (
    "net/http"
    "strings"
    "time"
)

type Request struct {
    Method          string
    StrURL          string

    Request         *http.Request
    Response        *http.Response
}

/**
 @Description：构造request
 @Param:
 @Return：
 */
func (r *Request) createRequest() error {
    var err error
    method := strings.ToUpper(r.Method)
    r.Request, err = http.NewRequest(method, r.StrURL, strings.NewReader("name=cjb"))
    if err != nil {
        return err
    }

    return nil
}

/**
 @Description：执行请求
 @Param:
 @Return：
 */
func (r *Request) Run() error {
    client := &http.Client{
        Timeout: time.Second * 5,                   // 设置5秒超时
    }

    err := r.createRequest()
    if err != nil {
        return nil
    }

    resp, err := client.Do(r.Request)
    if err != nil {
        return err
    }

    r.Response = resp
    return nil
}